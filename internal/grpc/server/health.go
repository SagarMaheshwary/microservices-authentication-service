package server

import (
	"context"

	user "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	healthpb.HealthServer
	UserClient  user.UserService
	RedisClient redis.RedisService
}

func (h *HealthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	status := healthpb.HealthCheckResponse_SERVING
	if err := h.RedisClient.Health(ctx); err != nil {
		status = healthpb.HealthCheckResponse_NOT_SERVING
	}

	if err := h.UserClient.Health(ctx); err != nil {
		status = healthpb.HealthCheckResponse_NOT_SERVING
	}

	logger.Info("Overall health status: %q", status)

	response := &healthpb.HealthCheckResponse{
		Status: status,
	}

	if status == healthpb.HealthCheckResponse_NOT_SERVING {
		prometheus.ServiceHealth.Set(0)
		return response, nil
	}

	prometheus.ServiceHealth.Set(1)
	return response, nil
}
