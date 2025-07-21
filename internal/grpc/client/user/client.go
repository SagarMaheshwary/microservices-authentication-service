package user

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func InitClient(ctx context.Context) (*grpc.ClientConn, error) {
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

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		logger.Error("User gRPC failed to connect on %q: %v", address, err)
		return nil, err
	}

	User = &userClient{
		client: pb.NewUserServiceClient(conn),
		health: healthpb.NewHealthClient(conn),
	}

	if err := HealthCheck(ctx); err != nil {
		return nil, err
	}

	logger.Info("User gRPC client connected on %q", address)

	return conn, err
}

func HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := User.health.Check(ctx, &healthpb.HealthCheckRequest{})

	if err != nil {
		logger.Error("User gRPC health check failed! %v", err)
		return err
	}

	if response.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("User gRPC health check failed!")
		return errors.New("User gRPC health check failed")
	}

	return nil
}
