package redis

import (
	"context"
	"fmt"
	"time"

	rds "github.com/redis/go-redis/v9"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
)

var ctx = context.Background()
var client *rds.Client

func Connect() {
	addr := fmt.Sprintf("%s:%d", config.GetRedis().Host, config.GetRedis().Port)

	client = rds.NewClient(&rds.Options{
		Addr:     addr,
		Username: config.GetRedis().Username,
		Password: config.GetRedis().Password,
	})

	if pong := client.Ping(ctx); pong.Val() != "PONG" {
		log.Fatal("Unable to connect redis %v", pong)
	}

	log.Info("Connected to redis on %q", addr)
}

func Get(key string) (string, error) {
	r, err := client.Get(ctx, key).Result()

	if err != nil {
		log.Error("Redis get key failed %v", err)
	}

	return r, err
}

func Set(key string, val string, exp time.Duration) error {
	err := client.Set(ctx, key, val, exp).Err()

	if err != nil {
		log.Error("Redis set key failed %v", err)
	}

	return err
}

func Del(key string) error {
	err := client.Del(ctx, key).Err()

	if err != nil {
		log.Error("Redis delete key failed %v", err)
	}

	return err
}
