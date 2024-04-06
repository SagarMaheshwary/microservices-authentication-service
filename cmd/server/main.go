package main

import (
	"fmt"
	"log"
	"net"
	"path"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helpers"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	loadenv()

	redis.Connect()

	client.InitClient()

	grpcServerConfig := config.GetgrpcServer()

	address := fmt.Sprintf("%v:%d", grpcServerConfig.Host, grpcServerConfig.Port)

	initgrpcServer(newListener(address))
}

func loadenv() {
	envPath := path.Join(helpers.RootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatalf("Unable to load env from %q: %v", envPath, err)
	}

	config.InitConfig()
}

func newListener(address string) *net.Listener {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to create tcp listner on %q: %v", address, err)
	}

	log.Printf("Starting tcp listener on %q", address)

	return &listener
}

func initgrpcServer(listener *net.Listener) {
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, server.NewAuthServer())

	log.Printf("grpc server started.")

	grpcServer.Serve(*listener)
}
