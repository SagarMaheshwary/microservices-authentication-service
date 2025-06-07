package user

import (
	"context"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect(ctx context.Context) {
	var opts []grpc.DialOption

	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)

	address := config.Conf.GRPCClient.UserServiceURL

	connection, err := grpc.Dial(address, opts...)

	if err != nil {
		logger.Error("User gRPC failed to connect on %q: %v", address, err)

		return
	}

	User = &userClient{
		client: pb.NewUserServiceClient(connection),
		health: healthpb.NewHealthClient(connection),
	}

	if HealthCheck(ctx) {
		logger.Info("User gRPC client connected on %q", address)
	}
}

func HealthCheck(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.Timeout)
	defer cancel()

	response, err := User.health.Check(ctx, &healthpb.HealthCheckRequest{})

	if err != nil {
		logger.Error("User gRPC health check failed! %v", err)

		return false
	}

	if response.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("User gRPC health check failed!")

		return false
	}

	return true
}
