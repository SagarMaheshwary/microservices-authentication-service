package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	user "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
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

var existingUserEmail = "taken@gmail.com"

func TestUserClient_FindById(t *testing.T) {
	req := &userpb.FindByIdRequest{Id: dummyUser.Id}
	res := &userpb.FindByIdResponse{
		Message: constant.MessageOK,
		Data:    &userpb.FindByIdResponseData{User: dummyUser},
	}

	cfg := &config.GRPCUserClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *userpb.FindByIdResponse
		mockErr    error
		expectErr  bool
		expectGRPC codes.Code
	}{
		{
			name:       "success",
			mockReturn: res,
			mockErr:    nil,
			expectErr:  false,
		},
		{
			name:       "gRPC error",
			mockReturn: nil,
			mockErr:    errors.New("grpc error"),
			expectErr:  true,
		},
		{
			name:       "not found",
			mockReturn: nil,
			mockErr:    status.Error(codes.NotFound, "user not found"),
			expectErr:  true,
			expectGRPC: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("FindById", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := user.NewUserClient(mockClient, mockHealth, cfg)

			got, err := c.FindById(context.Background(), req)

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUserClient_FindByCredential(t *testing.T) {
	req := &userpb.FindByCredentialRequest{
		Email:    dummyUser.Email,
		Password: "password",
	}
	res := &userpb.FindByCredentialResponse{
		Message: constant.MessageOK,
		Data: &userpb.FindByCredentialResponseData{
			User: dummyUser,
		},
	}

	cfg := &config.GRPCUserClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *userpb.FindByCredentialResponse
		mockErr    error
		expectErr  bool
		expectGRPC codes.Code
	}{
		{
			name:       "success",
			mockReturn: res,
			mockErr:    nil,
			expectErr:  false,
		},
		{
			name:       "gRPC error",
			mockReturn: nil,
			mockErr:    errors.New("grpc error"),
			expectErr:  true,
		},
		{
			name:       "not found",
			mockReturn: nil,
			mockErr:    status.Error(codes.NotFound, "user not found"),
			expectErr:  true,
			expectGRPC: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("FindByCredential", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := user.NewUserClient(mockClient, mockHealth, cfg)

			got, err := c.FindByCredential(context.Background(), req)

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, got)

				// if we expect a specific gRPC code
				if tt.expectGRPC != 0 {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.expectGRPC, st.Code())
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUserClient_Store(t *testing.T) {
	req := &userpb.StoreRequest{
		Name:     dummyUser.Name,
		Email:    dummyUser.Email,
		Password: "password",
	}
	res := &userpb.StoreResponse{
		Message: constant.MessageOK,
		Data: &userpb.StoreResponseData{
			User: dummyUser,
		},
	}

	cfg := &config.GRPCUserClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *userpb.StoreResponse
		mockErr    error
		expectErr  bool
		expectGRPC codes.Code
		req        *userpb.StoreRequest
	}{
		{
			name:       "success",
			mockReturn: res,
			mockErr:    nil,
			expectErr:  false,
			req:        req,
		},
		{
			name:       "gRPC error",
			mockReturn: nil,
			mockErr:    errors.New("grpc error"),
			expectErr:  true,
			req:        req,
		},
		{
			name:       "user exists",
			mockReturn: nil,
			mockErr:    status.Error(codes.NotFound, "user exists"),
			expectErr:  true,
			expectGRPC: codes.AlreadyExists,
			req: &userpb.StoreRequest{
				Name:     dummyUser.Name,
				Email:    existingUserEmail,
				Password: "password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("Store", mock.Anything, tt.req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := user.NewUserClient(mockClient, mockHealth, cfg)

			got, err := c.Store(context.Background(), tt.req)

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, got)

				// if we expect a specific gRPC code
				if tt.expectGRPC != 0 {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.expectGRPC, st.Code())
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUploadClient_Health(t *testing.T) {
	cfg := &config.GRPCUserClient{Timeout: 2 * time.Second}

	tests := []struct {
		name      string
		mockResp  *healthpb.HealthCheckResponse
		mockErr   error
		expectErr bool
		expectMsg string
	}{
		{
			name:      "success",
			mockResp:  &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING},
			mockErr:   nil,
			expectErr: false,
		},
		{
			name:      "gRPC error",
			mockResp:  nil,
			mockErr:   errors.New("grpc health check failed"),
			expectErr: true,
			expectMsg: "grpc health check failed",
		},
		{
			name:      "not serving",
			mockResp:  &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING},
			mockErr:   nil,
			expectErr: true,
			expectMsg: "user grpc health check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpload := new(MockUserServiceClient)
			mockHealth := new(MockHealthClient)

			client := user.NewUserClient(mockUpload, mockHealth, cfg)

			mockHealth.On("Check", mock.Anything, &healthpb.HealthCheckRequest{}).
				Return(tt.mockResp, tt.mockErr).Once()

			err := client.Health(context.Background())

			if tt.expectErr {
				require.Error(t, err)
				if tt.expectMsg != "" {
					assert.EqualError(t, err, tt.expectMsg)
				}
			} else {
				require.NoError(t, err)
			}

			mockHealth.AssertExpectations(t)
		})
	}
}
