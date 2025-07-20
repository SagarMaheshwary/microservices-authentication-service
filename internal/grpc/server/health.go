package server

import (
	"context"

	userrpc "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthServer struct {
	healthpb.HealthServer
}

func (h *healthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	status := getServicesHealthStatus(ctx)

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

func getServicesHealthStatus(ctx context.Context) healthpb.HealthCheckResponse_ServingStatus {
	if err := redis.HealthCheck(); err != nil {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	if err := userrpc.HealthCheck(ctx); err != nil {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	return healthpb.HealthCheckResponse_SERVING
}
