package main

import (
	"path"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helpers"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
)

func main() {
	log.Init()

	loadenv()

	redis.Connect()
	client.ConnectUserClient()
	server.Connect()
}

func loadenv() {
	envPath := path.Join(helpers.RootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatal("Failed to load .env %q: %v", envPath, err)
	}

	config.InitConfig()
}
