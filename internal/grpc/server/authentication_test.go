package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	libjwt "github.com/golang-jwt/jwt/v5"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"

	authpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/authentication"
	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
)

var now = time.Now().String()

var dummyUser = &userpb.User{
	Id:        1,
	Name:      "name",
	Email:     "name@gmail.com",
	Image:     nil,
	CreatedAt: &now,
	UpdatedAt: nil,
}

func TestAuthenticationServer_Register(t *testing.T) {
	mockResp := &userpb.StoreResponse{Data: &userpb.StoreResponseData{User: dummyUser}}

	tests := []struct {
		name          string
		setupMocks    func(u *MockUserClient, j *MockJWTManager)
		input         *authpb.RegisterRequest
		expectErr     bool
		expectedMsg   string
		expectedToken string
	}{
		{
			name: "success",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				u.On("Store", mock.Anything, mock.Anything).Return(mockResp, nil)
				j.On("NewToken", uint(dummyUser.Id), dummyUser.Email).Return("token123", nil)
			},
			input:         &authpb.RegisterRequest{Name: dummyUser.Name, Email: dummyUser.Email, Password: "pass"},
			expectErr:     false,
			expectedMsg:   constant.MessageOK,
			expectedToken: "token123",
		},
		{
			name: "user store fails",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				u.On("Store", mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
			},
			input:     &authpb.RegisterRequest{Name: dummyUser.Name, Email: dummyUser.Email, Password: "pass"},
			expectErr: true,
		},
		{
			name: "token generation fails",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				u.On("Store", mock.Anything, mock.Anything).Return(mockResp, nil)
				j.On("NewToken", uint(dummyUser.Id), dummyUser.Email).Return("", errors.New("jwt fail"))
			},
			input:     &authpb.RegisterRequest{Name: dummyUser.Name, Email: dummyUser.Email, Password: "pass"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := new(MockUserClient)
			j := new(MockJWTManager)
			if tt.setupMocks != nil {
				tt.setupMocks(u, j)
			}

			s := &server.AuthenticationServer{UserClient: u, JWTManager: j}
			resp, err := s.Register(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, resp.Message)
				assert.Equal(t, tt.expectedToken, resp.Data.Token)
			}

			u.AssertExpectations(t)
			j.AssertExpectations(t)
		})
	}
}

func TestAuthenticationServer_Login(t *testing.T) {
	mockResp := &userpb.FindByCredentialResponse{Data: &userpb.FindByCredentialResponseData{User: dummyUser}}

	tests := []struct {
		name          string
		setupMocks    func(u *MockUserClient, j *MockJWTManager)
		input         *authpb.LoginRequest
		expectErr     bool
		expectedMsg   string
		expectedToken string
	}{
		{
			name: "success",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				u.On("FindByCredential", mock.Anything, mock.Anything).Return(mockResp, nil)
				j.On("NewToken", uint(dummyUser.Id), dummyUser.Name).Return("token123", nil)
			},
			input:         &authpb.LoginRequest{Email: dummyUser.Email, Password: "pass"},
			expectErr:     false,
			expectedMsg:   constant.MessageOK,
			expectedToken: "token123",
		},
		{
			name: "user not found",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				u.On("FindByCredential", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
			},
			input:     &authpb.LoginRequest{Email: dummyUser.Email, Password: "pass"},
			expectErr: true,
		},
		{
			name: "token generation fails",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				u.On("FindByCredential", mock.Anything, mock.Anything).Return(mockResp, nil)
				j.On("NewToken", uint(dummyUser.Id), dummyUser.Name).Return("", errors.New("jwt fail"))
			},
			input:     &authpb.LoginRequest{Email: dummyUser.Email, Password: "pass"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := new(MockUserClient)
			j := new(MockJWTManager)
			if tt.setupMocks != nil {
				tt.setupMocks(u, j)
			}

			s := &server.AuthenticationServer{UserClient: u, JWTManager: j}
			resp, err := s.Login(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, resp.Message)
				assert.Equal(t, tt.expectedToken, resp.Data.Token)
			}

			u.AssertExpectations(t)
			j.AssertExpectations(t)
		})
	}
}

func TestAuthenticationServer_VerifyToken(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer validtoken"))

	tests := []struct {
		name        string
		setupMocks  func(u *MockUserClient, j *MockJWTManager)
		inputCtx    context.Context
		expectErr   bool
		expectedMsg string
	}{
		{
			name: "success",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				j.On("ParseToken", "validtoken").Return(libjwt.MapClaims{"id": float64(dummyUser.Id), "jti": "1"}, nil)
				j.On("IsBlacklisted", mock.Anything, "1").Return(false)
				u.On("FindById", mock.Anything, mock.Anything).Return(&userpb.FindByIdResponse{Data: &userpb.FindByIdResponseData{User: dummyUser}}, nil)
			},
			inputCtx:    ctx,
			expectErr:   false,
			expectedMsg: constant.MessageOK,
		},
		{
			name: "blacklisted token",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				j.On("ParseToken", "validtoken").Return(libjwt.MapClaims{"id": float64(dummyUser.Id), "jti": "1"}, nil)
				j.On("IsBlacklisted", mock.Anything, "1").Return(true)
			},
			inputCtx:  ctx,
			expectErr: true,
		},
		{
			name: "invalid token",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				j.On("ParseToken", "validtoken").Return(nil, errors.New("fail"))
			},
			inputCtx:  ctx,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := new(MockUserClient)
			j := new(MockJWTManager)
			if tt.setupMocks != nil {
				tt.setupMocks(u, j)
			}

			s := &server.AuthenticationServer{UserClient: u, JWTManager: j}
			resp, err := s.VerifyToken(tt.inputCtx, &authpb.VerifyTokenRequest{})

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, resp.Message)
			}

			u.AssertExpectations(t)
			j.AssertExpectations(t)
		})
	}
}

func TestAuthenticationServer_Logout(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer validtoken"))

	tests := []struct {
		name        string
		setupMocks  func(u *MockUserClient, j *MockJWTManager)
		inputCtx    context.Context
		expectErr   bool
		expectedMsg string
	}{
		{
			name: "success",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				j.On("ParseToken", "validtoken").Return(libjwt.MapClaims{"jti": "1", "exp": float64(time.Now().Add(time.Hour).Unix())}, nil)
				j.On("IsBlacklisted", mock.Anything, "1").Return(false)
				j.On("AddToBlacklist", mock.Anything, "1", mock.Anything).Return(nil)
			},
			inputCtx:    ctx,
			expectErr:   false,
			expectedMsg: constant.MessageOK,
		},
		{
			name: "blacklisted token",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				j.On("ParseToken", "validtoken").Return(libjwt.MapClaims{"jti": "1", "exp": float64(time.Now().Add(time.Hour).Unix())}, nil)
				j.On("IsBlacklisted", mock.Anything, "1").Return(true)
			},
			inputCtx:  ctx,
			expectErr: true,
		},
		{
			name: "token parse fails",
			setupMocks: func(u *MockUserClient, j *MockJWTManager) {
				j.On("ParseToken", "validtoken").Return(nil, errors.New("fail"))
			},
			inputCtx:  ctx,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := new(MockUserClient)
			j := new(MockJWTManager)
			if tt.setupMocks != nil {
				tt.setupMocks(u, j)
			}

			s := &server.AuthenticationServer{UserClient: u, JWTManager: j}
			resp, err := s.Logout(tt.inputCtx, &authpb.LogoutRequest{})

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, resp.Message)
			}

			u.AssertExpectations(t)
			j.AssertExpectations(t)
		})
	}
}
