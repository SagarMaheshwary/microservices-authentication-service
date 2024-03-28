package server

import (
	"context"
	"log"
	"strings"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/pkg"
	authpb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
	usrpb "github.com/sagarmaheshwary/microservices-authentication-service/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

	user := res.Data.User

	token, err := pkg.Createjwt(uint(user.Id), user.Name)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	loginResponse := &authpb.LoginResponse{
		Message: "Success",
		Data: &authpb.LoginResponseData{
			Token: token,
			User: &authpb.User{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				Image:     user.Image,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	}

	return loginResponse, nil
}

func (a *authServer) VerifyToken(ctx context.Context, data *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	bearerToken := md.Get("authorization")

	if len(bearerToken) == 0 {
		log.Println("Token is invalid.")

		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	log.Println(bearerToken[0])

	token, _ := strings.CutPrefix(bearerToken[0], "Bearer ")
	log.Println(token)

	claims, err := pkg.Parsejwt(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	userId := claims["id"].(float64)

	res, err := client.Client.FindById(&usrpb.FindByIdRequest{
		Id: int32(userId),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	user := res.Data.User

	return &authpb.VerifyTokenResponse{
		Message: "Success",
		Data: &authpb.VerifyTokenResponseData{
			User: &authpb.User{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				Image:     user.Image,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	}, nil
}

func (a *authServer) Logout(ctx context.Context, data *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	//@TODO: implement token blacklist

	return &authpb.LogoutResponse{
		Message: "Success",
		Data:    &authpb.LogoutResponseData{},
	}, nil
}

func NewAuthServer() *authServer {
	return &authServer{}
}
