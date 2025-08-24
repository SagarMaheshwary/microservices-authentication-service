package jwt_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

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
