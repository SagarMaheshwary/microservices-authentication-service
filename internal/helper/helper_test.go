package helper_test

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetRootDir(t *testing.T) {
	// Derive the expected path using runtime.Caller in the test itself
	_, b, _, _ := runtime.Caller(0)
	expected := filepath.Dir(path.Join(path.Dir(b)))

	t.Run("should return root directory of project", func(t *testing.T) {
		root := helper.GetRootDir()
		assert.Equal(t, expected, root)
	})
}

func TestGetGRPCMetadataValue(t *testing.T) {
	tests := []struct {
		name    string
		md      metadata.MD
		key     string
		wantVal string
		wantOk  bool
	}{
		{
			name:    "key exists with single value",
			md:      metadata.Pairs("authorization", "token123"),
			key:     "authorization",
			wantVal: "token123",
			wantOk:  true,
		},
		{
			name:    "key exists with multiple values - first returned",
			md:      metadata.Pairs("trace-id", "first", "trace-id", "second"),
			key:     "trace-id",
			wantVal: "first",
			wantOk:  true,
		},
		{
			name:    "key does not exist",
			md:      metadata.Pairs("authorization", "token123"),
			key:     "nonexistent",
			wantVal: "",
			wantOk:  false,
		},
		{
			name:    "empty metadata",
			md:      metadata.MD{},
			key:     "any",
			wantVal: "",
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := helper.GetGRPCMetadataValue(tt.md, tt.key)
			if gotVal != tt.wantVal || gotOk != tt.wantOk {
				t.Errorf("GetGRPCMetadataValue(%v, %q) = (%q, %v), want (%q, %v)",
					tt.md, tt.key, gotVal, gotOk, tt.wantVal, tt.wantOk)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal string
		expected   string
	}{
		{
			name:       "env variable set",
			envKey:     "TEST_ENV_STRING",
			envValue:   "value123",
			defaultVal: "default",
			expected:   "value123",
		},
		{
			name:       "env variable not set, uses default",
			envKey:     "TEST_ENV_STRING_NOT_SET",
			envValue:   "",
			defaultVal: "default",
			expected:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}
			got := helper.GetEnv(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal int
		expected   int
	}{
		{
			name:       "valid int from env",
			envKey:     "TEST_ENV_INT",
			envValue:   "42",
			defaultVal: 10,
			expected:   42,
		},
		{
			name:       "invalid int from env, fallback to default",
			envKey:     "TEST_ENV_INT_INVALID",
			envValue:   "abc",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "env not set, fallback to default",
			envKey:     "TEST_ENV_INT_NOT_SET",
			envValue:   "",
			defaultVal: 10,
			expected:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}
			got := helper.GetEnvInt(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetEnvDurationSeconds(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal time.Duration
		expected   time.Duration
	}{
		{
			name:       "valid duration from env",
			envKey:     "TEST_ENV_DURATION",
			envValue:   "3",
			defaultVal: 5,
			expected:   3 * time.Second,
		},
		{
			name:       "invalid duration from env, fallback",
			envKey:     "TEST_ENV_DURATION_INVALID",
			envValue:   "abc",
			defaultVal: 5,
			expected:   5 * time.Second,
		},
		{
			name:       "env not set, fallback",
			envKey:     "TEST_ENV_DURATION_NOT_SET",
			envValue:   "",
			defaultVal: 5,
			expected:   5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}
			got := helper.GetEnvDurationSeconds(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, got)
		})
	}
}
