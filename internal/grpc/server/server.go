package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-authentication-service/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/authentication"
	"google.golang.org/grpc"
)

func Connect() {
	c := config.GetgrpcServer()

	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthenticationServiceServer(grpcServer, &authenticationServer{})

	log.Info("gRPC server started on %q", address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Error("gRPC server failed to start %v", err)
	}
}
