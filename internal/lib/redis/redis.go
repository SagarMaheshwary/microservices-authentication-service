package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	redislib "github.com/redis/go-redis/v9"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
)

type RedisService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string, expiry time.Duration) error
	Del(ctx context.Context, key string) error
	Health(ctx context.Context) error
	Close() error
}

type RedisClient struct {
	Client *redislib.Client
}

func NewClient(cfg *config.Redis) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client := redislib.NewClient(&redislib.Options{
		Addr:     addr,
		Username: cfg.Username,
		Password: cfg.Password,
	})

	r := &RedisClient{Client: client}

	if err := r.Health(context.Background()); err != nil {
		logger.Error("%v", err)
		return nil, err
	}

	logger.Info("Redis server connected on %s", addr)

	return r, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, val string, expiry time.Duration) error {
	return r.Client.Set(ctx, key, val, expiry).Err()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClient) Health(ctx context.Context) error {
	for i := 0; i < 5; i++ {
		if pong := r.Client.Ping(ctx); pong.Val() == "PONG" {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("redis health check failed after retries")
}

func (r *RedisClient) Close() error {
	if err := r.Client.Close(); err != nil {
		logger.Error("Redis client close error %v", err)
		return err
	}

	return nil
}
