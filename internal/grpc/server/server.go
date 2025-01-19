package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	authpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/authentication"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func Connect() {
	c := config.Conf.GRPCServer

	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		logger.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(prometheusUnaryInterceptor))

	authpb.RegisterAuthenticationServiceServer(server, &authenticationServer{})
	healthpb.RegisterHealthServer(server, &healthServer{})

	logger.Info("gRPC server started on %q", address)

	if err := server.Serve(listener); err != nil {
		logger.Error("gRPC server failed to start %v", err)
	}
}

func prometheusUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	response, err := handler(ctx, req)

	method := info.FullMethod
	statusCode := status.Code(err).String()
	prometheus.GRPCRequestCounter.WithLabelValues(method, statusCode).Inc()
	prometheus.GRPCRequestLatency.WithLabelValues(method).Observe(time.Since(start).Seconds())

	return response, err
}
