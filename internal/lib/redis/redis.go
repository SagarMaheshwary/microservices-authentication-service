package redis

import (
	"context"
	"fmt"
	"time"

	redislib "github.com/redis/go-redis/v9"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
)

var ctx = context.Background()
var client *redislib.Client

func Connect() {
	c := config.Conf.Redis

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	client = redislib.NewClient(&redislib.Options{
		Addr:     addr,
		Username: c.Username,
		Password: c.Password,
	})

	if !HealthCheck() {
		logger.Error("Unable to connect to redis")

		return
	}

	logger.Info("Redis server connected on %q", addr)
}

func HealthCheck() bool {
	if pong := client.Ping(ctx); pong.Val() != "PONG" {
		logger.Error("Redis health check failed! %q", pong.Val())

		return false
	}

	return true
}

func Get(key string) (string, error) {
	r, err := client.Get(ctx, key).Result()

	if err != nil {
		logger.Error("Redis get key failed %v", err)
	}

	return r, err
}

func Set(key string, val string, exp time.Duration) error {
	err := client.Set(ctx, key, val, exp).Err()

	if err != nil {
		logger.Error("Redis set key failed %v", err)
	}

	return err
}

func Del(key string) error {
	err := client.Del(ctx, key).Err()

	if err != nil {
		logger.Error("Redis delete key failed %v", err)
	}

	return err
}
