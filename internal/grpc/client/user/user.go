package user

import (
	"context"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var User *userClient

type userClient struct {
	client pb.UserServiceClient
	health healthpb.HealthClient
}

func (u *userClient) FindById(ctx context.Context, data *pb.FindByIdRequest) (*pb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := u.client.FindById(ctx, data)
	if err != nil {
		logger.Error("gRPC userClient.FindById request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.FindById response: %v", response)
	return response, nil
}

func (u *userClient) FindByCredential(ctx context.Context, data *pb.FindByCredentialRequest) (*pb.FindByCredentialResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := u.client.FindByCredential(ctx, data)
	if err != nil {
		logger.Error("gRPC userClient.FindByCredential request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.FindByCredential response: %v", response)
	return response, nil
}

func (u *userClient) Store(ctx context.Context, data *pb.StoreRequest) (*pb.StoreResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := u.client.Store(ctx, data)
	if err != nil {
		logger.Error("gRPC userClient.Store request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.Store response: %v", response)
	return response, nil
}
