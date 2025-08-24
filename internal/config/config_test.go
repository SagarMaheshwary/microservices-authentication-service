package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigWithOptions(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	// Write a temporary .env file
	err := os.WriteFile(envFile, []byte("GRPC_SERVER_URL=192.168.0.1:5005\n"), 0644)
	require.NoError(t, err, "failed to write temp env file")

	tests := []struct {
		name                  string
		envFilePath           string
		setupEnv              func()
		expectedGRPCServerURL string
	}{
		{
			name:                  "loads from .env file",
			envFilePath:           envFile,
			setupEnv:              func() {}, // no system env
			expectedGRPCServerURL: "192.168.0.1:5005",
		},
		{
			name:                  "falls back to defaults when no .env file",
			envFilePath:           "nonexistent.env",
			setupEnv:              func() { os.Clearenv() },
			expectedGRPCServerURL: "0.0.0.0:5001", // default
		},
		{
			name:        "overrides with system env vars",
			envFilePath: "",
			setupEnv: func() {
				os.Clearenv()
				os.Setenv("GRPC_SERVER_URL", "some-host:6000")
			},
			expectedGRPCServerURL: "some-host:6000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			cfg := config.NewConfigWithOptions(config.LoaderOptions{EnvPath: tt.envFilePath})

			assert.Equal(t, tt.expectedGRPCServerURL, cfg.GRPCServer.URL, "GRPC server URL mismatch")

			// sanity check: durations and other defaults
			assert.NotEmpty(t, cfg.GRPCUserClient.URL, "GRPCUserClient.URL should not be empty")
			assert.Greater(t, cfg.GRPCUserClient.Timeout, 0*time.Second, "GRPCUserClient.Timeout should be > 0")
		})
	}
}
