package server

import (
	"context"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/pkg"
	authpb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
	usrpb "github.com/sagarmaheshwary/microservices-authentication-service/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct {
	authpb.AuthServiceServer
}

func (a *authServer) Login(ctx context.Context, data *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	res, err := client.Client.FindByCredential(&usrpb.FindByCredentialRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return nil, err
	}

	token, err := pkg.CreateJwt(uint(res.Data.User.Id), res.Data.User.Name)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	loginResponse := &authpb.LoginResponse{
		Message: "Success",
		Data: &authpb.LoginResponseData{
			Token: token,
			User: &authpb.User{
				Id:        res.Data.User.Id,
				Name:      res.Data.User.Name,
				Email:     res.Data.User.Email,
				Image:     res.Data.User.Image,
				CreatedAt: res.Data.User.CreatedAt,
				UpdatedAt: res.Data.User.UpdatedAt,
			},
		},
	}

	return loginResponse, nil
}

func NewAuthServer() *authServer {
	return &authServer{}
}
