package interceptors_test

import (
	"context"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server/interceptors"
	myprom "github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPrometheusUnaryInterceptor(t *testing.T) {
	reg := prometheus.NewRegistry()
	myprom.RegisterMetrics(reg)

	tests := []struct {
		name         string
		handlerErr   error
		expectedCode string
	}{
		{"success", nil, "OK"},
		{"internal_error", status.Error(codes.Internal, "oops"), "Internal"},
		{"custom_error", errors.New("custom"), "Unknown"}, // non-gRPC error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset metrics before each test
			myprom.GRPCRequestCounter.Reset()
			myprom.GRPCRequestLatency.Reset()

			// Create a dummy handler that returns specified error
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "response", tt.handlerErr
			}

			info := &grpc.UnaryServerInfo{FullMethod: "/test.Method"}

			resp, err := interceptors.PrometheusUnaryInterceptor(context.Background(), "req", info, handler)

			assert.Equal(t, "response", resp)
			assert.Equal(t, tt.handlerErr, err)

			metrics, err := reg.Gather()
			assert.NoError(t, err)

			counterFound := false
			for _, m := range metrics {
				if m.GetName() == "grpc_requests_total" {
					for _, metric := range m.GetMetric() {
						methodLabel := metric.GetLabel()[0].GetValue()
						statusLabel := metric.GetLabel()[1].GetValue()
						if methodLabel == info.FullMethod && statusLabel == tt.expectedCode {
							counterFound = true
						}
					}
				}
			}
			assert.True(t, counterFound, "counter should be incremented with correct labels")

			latencyFound := false
			for _, m := range metrics {
				if m.GetName() == "grpc_request_duration_seconds" {
					for _, metric := range m.GetMetric() {
						methodLabel := metric.GetLabel()[0].GetValue()
						if methodLabel == info.FullMethod {
							assert.GreaterOrEqual(t, metric.GetHistogram().GetSampleCount(), uint64(1))
							latencyFound = true
						}
					}
				}
			}
			assert.True(t, latencyFound, "latency metric should be observed")
		})
	}
}
