package server

import (
	"context"

	userrpc "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthServer struct {
	healthpb.HealthServer
}

func (h *healthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	status := getServicesHealthStatus()

	log.Info("Overall health status: %q", status)

	return &healthpb.HealthCheckResponse{
		Status: status,
	}, nil
}

func getServicesHealthStatus() healthpb.HealthCheckResponse_ServingStatus {
	if !redis.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	if !userrpc.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	return healthpb.HealthCheckResponse_SERVING
}
