package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdownJaeger := jaeger.Init(ctx)

	promServer := prometheus.NewServer()
	go prometheus.Serve(promServer)

	if err := redis.InitClient(); err != nil {
		os.Exit(constant.ExitFailure)
	}

	userConn, err := userrpc.InitClient(ctx)
	if err != nil {
		logger.Error("Failed to init User client: %v", err)
		os.Exit(constant.ExitFailure)
	}
	defer userConn.Close()

	grpcServer := server.NewServer()
	go func() {
		if err := server.Serve(grpcServer); err != nil {
			os.Exit(constant.ExitFailure)
		}
	}()

	<-ctx.Done()

	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdownJaeger(shutdownCtx); err != nil {
		logger.Warn("failed to shutdown jaeger tracer: %v", err)
	}

	shutdownCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := promServer.Shutdown(shutdownCtx); err != nil {
		logger.Warn("Prometheus server shutdown error: %v", err)
	}

	grpcServer.GracefulStop()

	logger.Info("Shutdown complete")
}
