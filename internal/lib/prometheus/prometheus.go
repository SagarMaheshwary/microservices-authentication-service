package prometheus

import (
	"fmt"
	"net/http"

	prometheuslib "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
)

var (
	GRPCRequestCounter = prometheuslib.NewCounterVec(
		prometheuslib.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GRPCRequestLatency = prometheuslib.NewHistogramVec(
		prometheuslib.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) of gRPC requests",
			Buckets: prometheuslib.DefBuckets,
		},
		[]string{"method"},
	)

	ServiceHealth = prometheuslib.NewGauge(prometheuslib.GaugeOpts{
		Name: "service_health_status",
		Help: "Health status of the service: 1=Healthy, 0=Unhealthy",
	})
)

func Connect() {
	c := config.Conf.Prometheus

	prometheuslib.MustRegister(GRPCRequestCounter, GRPCRequestLatency, ServiceHealth)

	address := fmt.Sprintf("%s:%d", c.METRICS_HOST, c.METRICS_PORT)

	http.Handle("/metrics", promhttp.Handler())

	logger.Info("Prometheus metrics endpoint running on %s", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		logger.Error("Failed to create http server for prometheus! %err", err)
	}
}
