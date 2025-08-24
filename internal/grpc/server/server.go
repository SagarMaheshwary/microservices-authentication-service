// server.go
package server

import (
	"net"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/grpc/server/interceptors"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/redis"
	authpb "github.com/sagarmaheshwary/microservices-authentication-service/internal/proto/authentication"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func NewServer(
	userClient user.UserService,
	redisClient redis.RedisService,
	jwtManager jwt.JWTManager,
) *grpc.Server {
	// Create gRPC server with interceptors & tracing
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.PrometheusUnaryInterceptor),
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)

	authpb.RegisterAuthenticationServiceServer(s, &AuthenticationServer{
		UserClient: userClient,
		JWTManager: jwtManager,
	})

	healthpb.RegisterHealthServer(s, &HealthServer{
		UserClient:  userClient,
		RedisClient: redisClient,
	})

	return s
}

func ServeListener(listener net.Listener, server *grpc.Server) error {
	logger.Info("gRPC server started on %q", listener.Addr().String())
	if err := server.Serve(listener); err != nil {
		logger.Error("gRPC server failed: %v", err)
		return err
	}
	return nil
}

func Serve(url string, server *grpc.Server) error {
	listener, err := net.Listen("tcp", url)
	if err != nil {
		logger.Error("Failed to create tcp listener on %q: %v", url, err)
		return err
	}
	return ServeListener(listener, server)
}
