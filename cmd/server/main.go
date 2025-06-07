package main

import (
	"context"
	"log"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	userrpc "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	server "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

func main() {
	logger.Init()
	config.Init()

	ctx := context.Background()
	shutdown := jaeger.Init(ctx)

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown jaeger tracer: %v", err)
		}
	}()

	go func() {
		prometheus.Connect()
	}()

	redis.Connect()
	userrpc.Connect(ctx)

	server.Connect()
}
