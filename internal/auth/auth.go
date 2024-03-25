package auth

import (
	"context"
	"log"

	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
)

type authServer struct {
	pb.AuthServiceServer
}

func (a *authServer) Login(ctx context.Context, data *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Println("LOGIN DATA", data)

	response := &pb.LoginResponse{
		Message: "Success",
		Data: &pb.LoginResponseData{
			Token: "Jwt token",
			User: &pb.User{
				Id:    1,
				Name:  "Daniel",
				Email: "daniel@gmail.com",
			},
		},
	}

	return response, nil
}

func NewServer() *authServer {
	return &authServer{}
}
