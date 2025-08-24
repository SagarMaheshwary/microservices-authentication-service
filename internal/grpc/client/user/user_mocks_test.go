package user_test

import (
	"context"

	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) FindById(ctx context.Context, in *userpb.FindByIdRequest, opts ...grpc.CallOption) (*userpb.FindByIdResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*userpb.FindByIdResponse), nil
}

func (m *MockUserServiceClient) FindByCredential(ctx context.Context, in *userpb.FindByCredentialRequest, opts ...grpc.CallOption) (*userpb.FindByCredentialResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*userpb.FindByCredentialResponse), nil
}

func (m *MockUserServiceClient) Store(ctx context.Context, in *userpb.StoreRequest, opts ...grpc.CallOption) (*userpb.StoreResponse, error) {
	args := m.Called(ctx, in)

	if in.Email == existingUserEmail {
		return nil, status.Error(codes.AlreadyExists, "user exists")
	}

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*userpb.StoreResponse), nil
}

func (m *MockUserServiceClient) Health(ctx context.Context, opts ...grpc.CallOption) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

type MockHealthClient struct {
	mock.Mock
	healthpb.HealthClient
}

func (m *MockHealthClient) Check(ctx context.Context, in *healthpb.HealthCheckRequest, opts ...grpc.CallOption) (*healthpb.HealthCheckResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*healthpb.HealthCheckResponse), args.Error(1)
}
