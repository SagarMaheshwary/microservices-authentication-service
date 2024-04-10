package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
	"google.golang.org/grpc"
)

func Connect() {
	grpcServerConfig := config.GetgrpcServer()

	address := fmt.Sprintf("%v:%d", grpcServerConfig.Host, grpcServerConfig.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	log.Info("gRPC server started on %q", address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Error("gRPC server failed to start %v", err)
	}
}
