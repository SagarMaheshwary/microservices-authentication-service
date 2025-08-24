package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	user "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestInitClient(t *testing.T) {
	dialErr := errors.New("dial failed")
	healthErr := errors.New("health failed")

	tests := []struct {
		name       string
		dialFunc   func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
		mockHealth func(h *MockHealthClient)
		expectErr  error
		expectNil  bool
	}{
		{
			name: "success",
			dialFunc: func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				return &grpc.ClientConn{}, nil
			},
			mockHealth: func(h *MockHealthClient) {
				h.On("Check", mock.Anything, mock.Anything, mock.Anything).
					Return(&healthpb.HealthCheckResponse{
						Status: healthpb.HealthCheckResponse_SERVING,
					}, nil)
			},
			expectErr: nil,
			expectNil: false,
		},
		{
			name: "health fails",
			dialFunc: func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				return &grpc.ClientConn{}, nil
			},
			mockHealth: func(h *MockHealthClient) {
				h.On("Check", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, healthErr)
			},
			expectErr: healthErr,
			expectNil: true,
		},
		{
			name: "dial fails",
			dialFunc: func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				return nil, dialErr
			},
			mockHealth: nil, // should never be called
			expectErr:  dialErr,
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			mockHealth := new(MockHealthClient)

			opt := &user.InitClientOptions{
				Config: &config.GRPCUserClient{URL: "fake-url"},
				Dial:   tt.dialFunc,
				Factory: func(c userpb.UserServiceClient, h healthpb.HealthClient, cfg *config.GRPCUserClient) user.UserService {
					return user.NewUserClient(mockClient, mockHealth, cfg)
				},
			}

			if tt.mockHealth != nil {
				tt.mockHealth(mockHealth)
			}

			svc, conn, err := user.NewClient(context.Background(), opt)

			if tt.expectErr != nil {
				assert.EqualError(t, err, tt.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectNil {
				assert.Nil(t, svc)
				assert.Nil(t, conn)
			} else {
				assert.NotNil(t, svc)
				assert.NotNil(t, conn)
			}

			mockHealth.AssertExpectations(t)
		})
	}
}
