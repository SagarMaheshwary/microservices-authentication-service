// server_test.go
package server_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewServer_ServiceRegistration(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
	}{
		{"AuthenticationService_registered", "auth.AuthenticationService"},
		{"HealthService_registered", "grpc.health.v1.Health"},
	}

	s := server.NewServer(nil, nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := s.GetServiceInfo()
			_, ok := info[tt.serviceName]
			assert.True(t, ok, "%s should be registered", tt.serviceName)
		})
	}
}

func TestServeListener_Integration(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	mockUserClient := new(MockUserClient)
	mockRedis := new(MockRedisClient)
	mockJWT := new(MockJWTManager)

	//HealthCheck rpc calls redis, userClient "Health"
	mockRedis.On("Health", mock.Anything).Return(nil)
	mockUserClient.On("Health", mock.Anything).Return(nil)

	s := server.NewServer(mockUserClient, mockRedis, mockJWT)

	go func() {
		_ = server.ServeListener(lis, s)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	//Call HealthCheck RPC
	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	assert.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)

	s.GracefulStop()

	mockRedis.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}
