package server_test

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	myprom "github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthServer_Check(t *testing.T) {
	tests := []struct {
		name           string
		redisErr       error
		userErr        error
		expectedStatus healthpb.HealthCheckResponse_ServingStatus
		expectedMetric float64
	}{
		{
			name:           "all healthy",
			redisErr:       nil,
			userErr:        nil,
			expectedStatus: healthpb.HealthCheckResponse_SERVING,
			expectedMetric: 1,
		},
		{
			name:           "redis down",
			redisErr:       assert.AnError,
			userErr:        nil,
			expectedStatus: healthpb.HealthCheckResponse_NOT_SERVING,
			expectedMetric: 0,
		},
		{
			name:           "user down",
			redisErr:       nil,
			userErr:        assert.AnError,
			expectedStatus: healthpb.HealthCheckResponse_NOT_SERVING,
			expectedMetric: 0,
		},
		{
			name:           "both down",
			redisErr:       assert.AnError,
			userErr:        assert.AnError,
			expectedStatus: healthpb.HealthCheckResponse_NOT_SERVING,
			expectedMetric: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisMock := new(MockRedisClient)
			redisMock.On("Health", mock.Anything).Return(tt.redisErr)

			userMock := new(MockUserClient)
			userMock.On("Health", mock.Anything).Return(tt.userErr)

			reg := prometheus.NewRegistry()
			myprom.RegisterMetrics(reg)

			hs := &server.HealthServer{
				RedisClient: redisMock,
				UserClient:  userMock,
			}

			resp, err := hs.Check(context.Background(), nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.Status)

			metrics, err := reg.Gather()
			assert.NoError(t, err)

			found := false
			for _, m := range metrics {
				if m.GetName() == "service_health_status" {
					value := m.GetMetric()[0].GetGauge().GetValue()
					assert.Equal(t, tt.expectedMetric, value)
					found = true
				}
			}
			assert.True(t, found, "service_health_status metric should be registered")
		})
	}
}
