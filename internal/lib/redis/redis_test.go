package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const redisImage = "redis:7.2"

func setupRedis(t *testing.T) (*redis.RedisClient, func()) {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        redisImage,
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(10 * time.Second),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	host, _ := redisC.Host(ctx)
	port, _ := redisC.MappedPort(ctx, "6379")

	client, err := redis.NewClient(&config.Redis{
		Host: host,
		Port: port.Int(),
	})

	if err != nil {
		t.Fatalf("failed to create redis client: %v", err)
	}

	teardown := func() {
		_ = client.Close()
		_ = redisC.Terminate(ctx)
	}

	return client, teardown
}

func TestRedisClient(t *testing.T) {
	client, teardown := setupRedis(t)
	defer teardown()

	ctx := context.Background()

	tests := []struct {
		name      string
		action    func() error
		expectErr bool
	}{
		{
			name: "set key",
			action: func() error {
				return client.Set(ctx, "foo", "bar", time.Second*5)
			},
			expectErr: false,
		},
		{
			name: "get key",
			action: func() error {
				val, err := client.Get(ctx, "foo")
				if err != nil {
					return err
				}
				if val != "bar" {
					t.Errorf("expected value 'bar', got %s", val)
				}
				return nil
			},
			expectErr: false,
		},
		{
			name: "delete key",
			action: func() error {
				return client.Del(ctx, "foo")
			},
			expectErr: false,
		},
		{
			name: "get deleted key returns error",
			action: func() error {
				_, err := client.Get(ctx, "foo")
				return err
			},
			expectErr: true, // redis returns redis.Nil for missing keys
		},
		{
			name: "health check",
			action: func() error {
				return client.Health(ctx)
			},
			expectErr: false,
		},
		{
			name: "close client",
			action: func() error {
				return client.Close()
			},
			expectErr: false,
		},
		{
			name: "close already closed client",
			action: func() error {
				return client.Close()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action()
			if tt.expectErr {
				require.Error(t, err, "expected an error but got none")
			} else {
				require.NoError(t, err, "did not expect error but got one")
			}
		})
	}
}
