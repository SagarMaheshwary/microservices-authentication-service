package prometheus_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prometheus/client_golang/prometheus"
	myprom "github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/prometheus"
)

func TestRegisterMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()

	myprom.RegisterMetrics(reg)

	// Touch the vectors so they are instantiated
	myprom.GRPCRequestCounter.WithLabelValues("methodX", "success").Inc()
	myprom.GRPCRequestLatency.WithLabelValues("methodX").Observe(0.123)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	tests := []struct {
		name string
		want string
	}{
		{"grpc_requests_total", "grpc_requests_total"},
		{"grpc_request_duration_seconds", "grpc_request_duration_seconds"},
		{"service_health_status", "service_health_status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, m := range metrics {
				if m.GetName() == tt.want {
					found = true
					break
				}
			}
			assert.True(t, found, "metric %s should be registered", tt.want)
		})
	}
}

func TestNewServerAndMetricsEndpoint(t *testing.T) {
	reg := prometheus.NewRegistry()
	server := myprom.NewServer(":0", reg) // :0 = random free port

	// Start test server with httptest instead of real listen
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	// Increment some metrics
	myprom.GRPCRequestCounter.WithLabelValues("Login", "OK").Inc()
	myprom.ServiceHealth.Set(1)

	resp, err := http.Get(ts.URL + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "grpc_requests_total")
	assert.Contains(t, bodyStr, "service_health_status")
}

func TestServeFunction(t *testing.T) {
	t.Run("successful listen", func(t *testing.T) {
		server := &http.Server{Addr: ":0"}
		called := false
		err := myprom.Serve(server, func() error {
			called = true
			return nil
		})
		assert.NoError(t, err)
		assert.True(t, called, "listen function should be called")
	})

	t.Run("server closed error", func(t *testing.T) {
		server := &http.Server{Addr: ":0"}
		err := myprom.Serve(server, func() error {
			return http.ErrServerClosed
		})
		assert.NoError(t, err, "ErrServerClosed should be ignored")
	})

	t.Run("unexpected error", func(t *testing.T) {
		server := &http.Server{Addr: ":0"}
		expectedErr := assert.AnError
		err := myprom.Serve(server, func() error {
			return expectedErr
		})
		assert.ErrorIs(t, err, expectedErr)
	})
}
