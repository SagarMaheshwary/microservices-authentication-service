package main

import (
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	userrpc "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	grpcsrv "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

func main() {
	log.Init()
	config.Init()

	redis.Connect()
	userrpc.Connect()
	grpcsrv.Connect()
}
