package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/authentication"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect() {
	c := config.Conf.GRPCServer

	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		logger.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var opts []grpc.ServerOption

	server := grpc.NewServer(opts...)
	pb.RegisterAuthenticationServiceServer(server, &authenticationServer{})
	healthpb.RegisterHealthServer(server, &healthServer{})

	logger.Info("gRPC server started on %q", address)

	if err := server.Serve(listener); err != nil {
		logger.Error("gRPC server failed to start %v", err)
	}
}
