package config

import (
	"os"
	"path"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
)

type Config struct {
	GRPCServer     *GRPCServer
	JWT            *JWT
	GRPCUserClient *GRPCUserClient
	Redis          *Redis
	Prometheus     *Prometheus
	Jaeger         *Jaeger
}

type GRPCServer struct {
	URL string
}

type GRPCUserClient struct {
	URL     string
	Timeout time.Duration
}

type JWT struct {
	Secret string
	Expiry time.Duration
}

type Redis struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Prometheus struct {
	URL string
}

type Jaeger struct {
	URL string
}

type LoaderOptions struct {
	EnvPath     string
	EnvLoader   func(string) error
	FileChecker func(string) bool
}

func NewConfigWithOptions(opts LoaderOptions) *Config {
	envLoader := opts.EnvLoader
	if envLoader == nil {
		envLoader = func(path string) error { return env.Load(path) }
	}
	fileChecker := opts.FileChecker
	if fileChecker == nil {
		fileChecker = func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		}
	}

	if opts.EnvPath != "" && fileChecker(opts.EnvPath) {
		if err := envLoader(opts.EnvPath); err != nil {
			logger.Panic("Failed to load .env %q: %v", opts.EnvPath, err)
		}
		logger.Info("Loaded environment variables from %q", opts.EnvPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	return &Config{
		GRPCServer: &GRPCServer{
			URL: helper.GetEnv("GRPC_SERVER_URL", "0.0.0.0:5001"),
		},
		JWT: &JWT{
			Secret: helper.GetEnv("JWT_SECRET", "secret-key"),
			Expiry: helper.GetEnvDurationSeconds("JWT_EXPIRY_SECONDS", 3600),
		},
		GRPCUserClient: &GRPCUserClient{
			URL:     helper.GetEnv("GRPC_USER_SERVICE_URL", "user-service:5000"),
			Timeout: helper.GetEnvDurationSeconds("GRPC_USER_SERVICE_TIMEOUT_SECONDS", 5),
		},
		Redis: &Redis{
			Host:     helper.GetEnv("REDIS_HOST", "redis"),
			Port:     helper.GetEnvInt("REDIS_PORT", 6379),
			Username: helper.GetEnv("REDIS_USERNAME", "default"),
			Password: helper.GetEnv("REDIS_PASSWORD", "password"),
		},
		Prometheus: &Prometheus{
			URL: helper.GetEnv("PROMETHEUS_URL", "0.0.0.0:5011"),
		},
		Jaeger: &Jaeger{
			URL: helper.GetEnv("JAEGER_URL", "jaeger:4318"),
		},
	}
}

func NewConfig() *Config {
	return NewConfigWithOptions(LoaderOptions{
		EnvPath: path.Join(helper.GetRootDir(), "..", ".env"),
	})
}
