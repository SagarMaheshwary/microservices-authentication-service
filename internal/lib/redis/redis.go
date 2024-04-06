package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	rds "github.com/redis/go-redis/v9"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
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
		log.Fatalf("Unable to Connect Redis %v", pong)
	}

	log.Println("Connected to Redis.")
}

func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Set(key string, val string, exp time.Duration) error {
	return client.Set(ctx, key, val, exp).Err()
}

func Del(key string) error {
	return client.Del(ctx, key).Err()
}
