package main

import (
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	grpcct "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	grpcsrv "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

func main() {
	log.Init()
	config.Init()

	redis.Connect()
	grpcct.Connect()
	grpcsrv.Connect()
}
