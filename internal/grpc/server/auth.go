package server

import (
	"context"
	"strings"

	cons "github.com/sagarmaheshwary/microservices-authentication-service/internal/constants"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helpers"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
	apb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
	upb "github.com/sagarmaheshwary/microservices-authentication-service/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	REGISTER_RPC_TOKEN_ERROR = "User successfully registered, but there was a problem creating the authentication token. Please try manual login."
)

type authServer struct {
	apb.AuthServiceServer
}

func (a *authServer) Register(ctx context.Context, data *apb.RegisterRequest) (*apb.RegisterResponse, error) {
	clientResponse, err := client.Client.Store(&upb.StoreRequest{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User

	token, err := jwt.New(uint(user.Id), user.Email)

	if err != nil {
		return nil, status.Errorf(codes.Internal, REGISTER_RPC_TOKEN_ERROR)
	}

	response := &apb.RegisterResponse{
		Message: cons.OK,
		Data: &apb.RegisterResponseData{
			Token: token,
			User: &apb.User{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				Image:     user.Image,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	}

	return response, nil
}

func (a *authServer) Login(ctx context.Context, data *apb.LoginRequest) (*apb.LoginResponse, error) {
	clientResponse, err := client.Client.FindByCredential(&upb.FindByCredentialRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User

	token, err := jwt.New(uint(user.Id), user.Name)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.INTERNAL_SERVER_ERROR)
	}

	response := &apb.LoginResponse{
		Message: cons.OK,
		Data: &apb.LoginResponseData{
			Token: token,
			User: &apb.User{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				Image:     user.Image,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	}

	return response, nil
}

func (a *authServer) VerifyToken(ctx context.Context, data *apb.VerifyTokenRequest) (*apb.VerifyTokenResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	header, _ := helpers.GetFromMetadata(md, cons.HDR_AUTHORIZATION)
	token, f := strings.CutPrefix(header, cons.HDR_BEARER_PREFIX)

	if !f {
		return nil, status.Errorf(codes.Unauthenticated, cons.UNAUTHENTICATED)
	}

	claims, err := jwt.Parse(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.UNAUTHENTICATED)
	}

	if blacklisted := jwt.IsBlacklisted(claims["jti"].(string)); blacklisted {
		return nil, status.Errorf(codes.Unauthenticated, cons.UNAUTHENTICATED)
	}

	userId := claims["id"].(float64)

	clientResponse, err := client.Client.FindById(&upb.FindByIdRequest{
		Id: int32(userId),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.UNAUTHENTICATED)
	}

	user := clientResponse.Data.User

	response := &apb.VerifyTokenResponse{
		Message: cons.OK,
		Data: &apb.VerifyTokenResponseData{
			User: &apb.User{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				Image:     user.Image,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	}

	return response, nil
}

func (a *authServer) Logout(ctx context.Context, data *apb.LogoutRequest) (*apb.LogoutResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	header, _ := helpers.GetFromMetadata(md, cons.HDR_AUTHORIZATION)
	token, f := strings.CutPrefix(header, cons.HDR_BEARER_PREFIX)

	if !f {
		return nil, status.Errorf(codes.Unauthenticated, cons.UNAUTHENTICATED)
	}

	claims, err := jwt.Parse(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.UNAUTHENTICATED)
	}

	err = jwt.AddToBlacklist(claims["jti"].(string), int64(claims["exp"].(float64)))

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.INTERNAL_SERVER_ERROR)
	}

	response := &apb.LogoutResponse{
		Message: cons.OK,
		Data:    &apb.LogoutResponseData{},
	}

	return response, nil
}
