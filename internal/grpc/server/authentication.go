package server

import (
	"context"
	"strings"

	libjwt "github.com/golang-jwt/jwt/v5"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
	authpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/authentication"
	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	REGISTER_RPC_TOKEN_ERROR = "User successfully registered, but there was a problem creating the authentication token. Please try manual login."
)

type AuthenticationServer struct {
	authpb.AuthenticationServiceServer
	UserClient user.UserService
	JWTManager jwt.JWTManager
}

func (a *AuthenticationServer) Register(ctx context.Context, data *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	clientResponse, err := a.UserClient.Store(ctx, &userpb.StoreRequest{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User
	token, err := a.JWTManager.NewToken(uint(user.Id), user.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, REGISTER_RPC_TOKEN_ERROR)
	}

	response := &authpb.RegisterResponse{
		Message: constant.MessageOK,
		Data: &authpb.RegisterResponseData{
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
	return response, nil
}

func (a *AuthenticationServer) Login(ctx context.Context, data *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	clientResponse, err := a.UserClient.FindByCredential(ctx, &userpb.FindByCredentialRequest{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User
	token, err := a.JWTManager.NewToken(uint(user.Id), user.Name)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, constant.MessageUnauthorized)
	}

	response := &authpb.LoginResponse{
		Message: constant.MessageOK,
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
	return response, nil
}

func (a *AuthenticationServer) VerifyToken(ctx context.Context, data *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	claims, err := parseAndValidateJwtTokenFromMetadata(ctx, a.JWTManager)
	if err != nil {
		return nil, err
	}

	userId := claims["id"].(float64)

	clientResponse, err := a.UserClient.FindById(ctx, &userpb.FindByIdRequest{
		Id: int32(userId),
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, constant.MessageUnauthorized)
	}

	user := clientResponse.Data.User
	response := &authpb.VerifyTokenResponse{
		Message: constant.MessageOK,
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
	}
	return response, nil
}

func (a *AuthenticationServer) Logout(ctx context.Context, data *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	claims, err := parseAndValidateJwtTokenFromMetadata(ctx, a.JWTManager)
	if err != nil {
		return nil, err
	}

	err = a.JWTManager.AddToBlacklist(ctx, claims["jti"].(string), int64(claims["exp"].(float64)))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, constant.MessageUnauthorized)
	}

	response := &authpb.LogoutResponse{
		Message: constant.MessageOK,
		Data:    &authpb.LogoutResponseData{},
	}
	return response, nil
}

func parseAndValidateJwtTokenFromMetadata(ctx context.Context, jwtManager jwt.JWTManager) (libjwt.MapClaims, error) {
	authErr := status.Error(codes.Unauthenticated, constant.MessageUnauthorized)

	md, _ := metadata.FromIncomingContext(ctx)
	header, _ := helper.GetGRPCMetadataValue(md, constant.HeaderAuthorization)
	token, f := strings.CutPrefix(header, constant.HeaderBearerPrefix)
	if !f {
		return nil, authErr
	}

	claims, err := jwtManager.ParseToken(token)
	if err != nil {
		return nil, authErr
	}

	if blacklisted := jwtManager.IsBlacklisted(ctx, claims["jti"].(string)); blacklisted {
		return nil, authErr
	}

	return claims, nil
}
