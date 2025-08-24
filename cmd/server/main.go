package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	prometheuslib "github.com/prometheus/client_golang/prometheus"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	user "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	server "github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	"google.golang.org/grpc"
)

func main() {
	logger.Init()
	cfg := config.NewConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdownJaeger := jaeger.Init(ctx, cfg.Jaeger.URL)

	promServer := prometheus.NewServer(cfg.Prometheus.URL, prometheuslib.NewRegistry())
	go func() {
		if err := prometheus.Serve(promServer, promServer.ListenAndServe); err != nil && err != http.ErrServerClosed {
			stop()
		}
	}()

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		os.Exit(constant.ExitFailure)
	}
	defer redisClient.Close()

	userClient, userConn, err := user.NewClient(ctx, &user.InitClientOptions{Config: cfg.GRPCUserClient})
	if err != nil {
		logger.Error("Failed to connect to user client: %v", err)
		os.Exit(constant.ExitFailure)
	}
	defer userConn.Close()

	jwtManager := jwt.NewJWTManager(cfg.JWT, redisClient)

	grpcServer := server.NewServer(userClient, redisClient, jwtManager)
	go func() {
		if err := server.Serve(cfg.GRPCServer.URL, grpcServer); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			stop()
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
