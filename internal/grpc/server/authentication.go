package server

import (
	"context"
	"strings"

	cons "github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	userrpc "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
	apb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/authentication"
	upb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
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
	clientResponse, err := userrpc.User.Store(&upb.StoreRequest{
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
		Message: cons.MessageOK,
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
	clientResponse, err := userrpc.User.FindByCredential(&upb.FindByCredentialRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User

	token, err := jwt.New(uint(user.Id), user.Name)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageInternalServerError)
	}

	response := &apb.LoginResponse{
		Message: cons.MessageOK,
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

	header, _ := helper.GetFromMetadata(md, cons.HeaderAuthorization)
	token, f := strings.CutPrefix(header, cons.HeaderBearerPrefix)

	if !f {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageUnauthorized)
	}

	claims, err := jwt.Parse(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageUnauthorized)
	}

	if blacklisted := jwt.IsBlacklisted(claims["jti"].(string)); blacklisted {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageUnauthorized)
	}

	userId := claims["id"].(float64)

	clientResponse, err := userrpc.User.FindById(&upb.FindByIdRequest{
		Id: int32(userId),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MessageUnauthorized)
	}

	user := clientResponse.Data.User

	response := &apb.VerifyTokenResponse{
		Message: cons.MessageOK,
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

	header, _ := helper.GetFromMetadata(md, cons.HeaderAuthorization)
	token, f := strings.CutPrefix(header, cons.HeaderBearerPrefix)

	if !f {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageUnauthorized)
	}

	claims, err := jwt.Parse(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, cons.MessageUnauthorized)
	}

	err = jwt.AddToBlacklist(claims["jti"].(string), int64(claims["exp"].(float64)))

	if err != nil {
		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	response := &apb.LogoutResponse{
		Message: cons.MessageOK,
		Data:    &apb.LogoutResponseData{},
	}

	return response, nil
}
