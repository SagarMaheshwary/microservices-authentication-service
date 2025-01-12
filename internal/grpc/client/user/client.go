package user

import (
	"context"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect() {
	var options []grpc.DialOption

	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address := config.Conf.GRPCClient.UserServiceURL

	connection, err := grpc.Dial(address, options...)

	if err != nil {
		logger.Error("User gRPC failed to connect on %q: %v", address, err)
	}

	User = &userClient{
		client: pb.NewUserServiceClient(connection),
		health: healthpb.NewHealthClient(connection),
	}

	if HealthCheck() {
		logger.Info("User gRPC client connected on %q", address)
	}
}

func HealthCheck() bool {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)
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
