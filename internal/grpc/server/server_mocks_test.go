package server_test

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/user"
	"github.com/stretchr/testify/mock"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// ===== Mock User Client =====
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) FindById(ctx context.Context, in *userpb.FindByIdRequest) (*userpb.FindByIdResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*userpb.FindByIdResponse), nil
}

func (m *MockUserClient) FindByCredential(ctx context.Context, in *userpb.FindByCredentialRequest) (*userpb.FindByCredentialResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*userpb.FindByCredentialResponse), nil
}

func (m *MockUserClient) Store(ctx context.Context, in *userpb.StoreRequest) (*userpb.StoreResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*userpb.StoreResponse), nil
}

func (m *MockUserClient) Health(ctx context.Context) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

// ===== Mock Health Client =====
type MockHealthClient struct {
	mock.Mock
	healthpb.HealthClient
}

func (m *MockHealthClient) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*healthpb.HealthCheckResponse), args.Error(1)
}

// ===== Mock Redis Client =====
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)

	if err := args.Error(1); err != nil {
		return "", err
	}

	return args.Get(0).(string), nil
}

func (m *MockRedisClient) Set(ctx context.Context, key string, val string, expiration time.Duration) error {
	args := m.Called(ctx, key, val, expiration)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil

}

func (m *MockRedisClient) Del(ctx context.Context, keys string) error {
	args := m.Called(ctx, keys)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

func (m *MockRedisClient) Health(ctx context.Context) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

func (m *MockRedisClient) Close() error {
	return nil
}

// ===== Mock JWT Manager =====
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) NewToken(id uint, username string) (string, error) {
	args := m.Called(id, username)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) ParseToken(token string) (jwt.MapClaims, error) {
	args := m.Called(token)
	if claims, ok := args.Get(0).(jwt.MapClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJWTManager) AddToBlacklist(ctx context.Context, jti string, expiry int64) error {
	args := m.Called(ctx, jti, expiry)
	return args.Error(0)
}

func (m *MockJWTManager) IsBlacklisted(ctx context.Context, jti string) bool {
	args := m.Called(ctx, jti)
	return args.Bool(0)
}
