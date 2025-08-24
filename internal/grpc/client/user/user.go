package user

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type UserService interface {
	FindById(ctx context.Context, in *userpb.FindByIdRequest) (*userpb.FindByIdResponse, error)
	FindByCredential(ctx context.Context, in *userpb.FindByCredentialRequest) (*userpb.FindByCredentialResponse, error)
	Store(ctx context.Context, in *userpb.StoreRequest) (*userpb.StoreResponse, error)
	Health(ctx context.Context) error
}

type UserClient struct {
	config *config.GRPCUserClient
	client userpb.UserServiceClient
	health healthpb.HealthClient
}

func NewUserClient(c userpb.UserServiceClient, h healthpb.HealthClient, cfg *config.GRPCUserClient) *UserClient {
	return &UserClient{client: c, health: h, config: cfg}
}

func (u *UserClient) FindById(ctx context.Context, in *userpb.FindByIdRequest) (*userpb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	response, err := u.client.FindById(ctx, in)
	if err != nil {
		logger.Error("gRPC userClient.FindById request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.FindById response: %v", response)
	return response, nil
}

func (u *UserClient) FindByCredential(ctx context.Context, in *userpb.FindByCredentialRequest) (*userpb.FindByCredentialResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	response, err := u.client.FindByCredential(ctx, in)
	if err != nil {
		logger.Error("gRPC userClient.FindByCredential request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.FindByCredential response: %v", response)
	return response, nil
}

func (u *UserClient) Store(ctx context.Context, in *userpb.StoreRequest) (*userpb.StoreResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	response, err := u.client.Store(ctx, in)
	if err != nil {
		logger.Error("gRPC userClient.Store request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.Store response: %v", response)
	return response, nil
}

func (u *UserClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	res, err := u.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		logger.Error("User gRPC health check failed! %v", err)
		return err
	}

	if res.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("User gRPC health check failed")
		return errors.New("user grpc health check failed")
	}

	return nil
}
