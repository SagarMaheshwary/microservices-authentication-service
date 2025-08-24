package prometheus

import (
	"net/http"

	prometheuslib "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func RegisterMetrics(registry *prometheuslib.Registry) {
	registry.MustRegister(
		GRPCRequestCounter,
		GRPCRequestLatency,
		ServiceHealth,
	)
}

func NewServer(url string, registry *prometheuslib.Registry) *http.Server {
	RegisterMetrics(registry)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{},
	))

	return &http.Server{Addr: url, Handler: mux}
}

func Serve(server *http.Server, listen func() error) error {
	logger.Info("Starting Prometheus metrics server %s", server.Addr)

	if err := listen(); err != nil && err != http.ErrServerClosed {
		logger.Error("Prometheus http server error! %v", err)
		return err
	}
	return nil
}
