package server

import (
	"context"
	"strings"

	cons "github.com/sagarmaheshwary/microservices-authentication-service/internal/constants"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
	apb "github.com/sagarmaheshwary/microservices-authentication-service/proto/authentication"
	upb "github.com/sagarmaheshwary/microservices-authentication-service/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	REGISTER_RPC_TOKEN_ERROR = "User successfully registered, but there was a problem creating the authentication token. Please try manual login."
)

type authenticationServer struct {
	apb.AuthenticationServiceServer
}

func (a *authenticationServer) Register(ctx context.Context, data *apb.RegisterRequest) (*apb.RegisterResponse, error) {
	clientResponse, err := client.User.Store(&upb.StoreRequest{
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
		Message: cons.MSGOK,
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

func (a *authenticationServer) Login(ctx context.Context, data *apb.LoginRequest) (*apb.LoginResponse, error) {
	clientResponse, err := client.User.FindByCredential(&upb.FindByCredentialRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User

	token, err := jwt.New(uint(user.Id), user.Name)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.MSGInternalServerError)
	}

	response := &apb.LoginResponse{
		Message: cons.MSGOK,
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

func (a *authenticationServer) VerifyToken(ctx context.Context, data *apb.VerifyTokenRequest) (*apb.VerifyTokenResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	header, _ := helper.GetFromMetadata(md, cons.HDR_AUTHORIZATION)
	token, f := strings.CutPrefix(header, cons.HDR_BEARER_PREFIX)

	if !f {
		return nil, status.Errorf(codes.Unauthenticated, cons.MSGUnauthenticated)
	}

	claims, err := jwt.Parse(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.MSGUnauthenticated)
	}

	if blacklisted := jwt.IsBlacklisted(claims["jti"].(string)); blacklisted {
		return nil, status.Errorf(codes.Unauthenticated, cons.MSGUnauthenticated)
	}

	userId := claims["id"].(float64)

	clientResponse, err := client.User.FindById(&upb.FindByIdRequest{
		Id: int32(userId),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MSGUnauthenticated)
	}

	user := clientResponse.Data.User

	response := &apb.VerifyTokenResponse{
		Message: cons.MSGOK,
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

func (a *authenticationServer) Logout(ctx context.Context, data *apb.LogoutRequest) (*apb.LogoutResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	header, _ := helper.GetFromMetadata(md, cons.HDR_AUTHORIZATION)
	token, f := strings.CutPrefix(header, cons.HDR_BEARER_PREFIX)

	if !f {
		return nil, status.Errorf(codes.Unauthenticated, cons.MSGUnauthenticated)
	}

	claims, err := jwt.Parse(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.MSGUnauthenticated)
	}

	err = jwt.AddToBlacklist(claims["jti"].(string), int64(claims["exp"].(float64)))

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MSGInternalServerError)
	}

	response := &apb.LogoutResponse{
		Message: cons.MSGOK,
		Data:    &apb.LogoutResponseData{},
	}

	return response, nil
}
