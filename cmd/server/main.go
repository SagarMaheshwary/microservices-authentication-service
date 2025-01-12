package main

import (
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	userrpc "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	server "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

func main() {
	logger.Init()
	config.Init()

	redis.Connect()
	userrpc.Connect()
	server.Connect()
}
