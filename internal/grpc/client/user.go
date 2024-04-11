package client

import (
	"context"

	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/user"
)

var User *userClient

type userClient struct {
	client pb.UserServiceClient
}

func (u *userClient) FindById(data *pb.FindByIdRequest) (*pb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := u.client.FindById(ctx, data)

	if err != nil {
		log.Error("gRPC userClient.FindById request failed: %v", err)
		return nil, err
	}

	log.Info("gRPC userClient.FindById response: %v", response)

	return response, nil
}

func (u *userClient) FindByCredential(data *pb.FindByCredentialRequest) (*pb.FindByCredentialResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := u.client.FindByCredential(ctx, data)

	if err != nil {
		log.Error("gRPC userClient.FindByCredential request failed: %v", err)
		return nil, err
	}

	log.Info("gRPC userClient.FindByCredential response: %v", response)

	return response, nil
}

func (u *userClient) Store(data *pb.StoreRequest) (*pb.StoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := u.client.Store(ctx, data)

	if err != nil {
		log.Error("gRPC userClient.Store request failed: %v", err)
		return nil, err
	}

	log.Info("gRPC userClient.Store response: %v", response)

	return response, nil
}
