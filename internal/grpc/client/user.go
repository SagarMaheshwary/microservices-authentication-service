package client

import (
	"context"
	"log"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client *userClient

type userClient struct {
	client pb.UserServiceClient
}

func (u *userClient) FindById(data *pb.FindByIdRequest) (*pb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := u.client.FindById(ctx, data)

	if err != nil {
		log.Printf("userClient.FindById failed: %v", err)
		return nil, err
	}

	log.Printf("userClient.FindById response: %v", response)

	return response, nil
}

func (u *userClient) FindByCredential(data *pb.FindByCredentialRequest) (*pb.FindByCredentialResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := u.client.FindByCredential(ctx, data)

	if err != nil {
		log.Printf("userClient.FindByCredential failed: %v", err)
		return nil, err
	}

	log.Printf("userClient.FindByCredential response: %v", response)

	return response, nil
}

func InitClient() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address := config.GetgrpcClient().UserServiceurl

	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		log.Printf("grpc client connection failed on %q: %v", address, err)
	}

	log.Printf("Connected to grpc client: %q", address)

	Client = &userClient{client: pb.NewUserServiceClient(conn)}
}
