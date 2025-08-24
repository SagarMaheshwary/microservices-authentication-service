package interceptors

import (
	"context"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func PrometheusUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	response, err := handler(ctx, req)

	method := info.FullMethod
	statusCode := status.Code(err).String()
	prometheus.GRPCRequestCounter.WithLabelValues(method, statusCode).Inc()
	prometheus.GRPCRequestLatency.WithLabelValues(method).Observe(time.Since(start).Seconds())

	return response, err
}
