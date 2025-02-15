package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/logger"
)

var Conf *Config

type Config struct {
	GRPCServer *GRPCServer
	JWT        *JWT
	GRPCClient *GRPCClient
	Redis      *Redis
	Prometheus *Prometheus
}

type GRPCServer struct {
	Host string
	Port int
}

type GRPCClient struct {
	UserServiceURL string
	Timeout        time.Duration
}

type JWT struct {
	Secret string
	Expiry int
}

type Redis struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Prometheus struct {
	METRICS_HOST string
	METRICS_PORT int
}

func Init() {
	envPath := path.Join(helper.GetRootDir(), "..", ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := env.Load(envPath); err != nil {
			logger.Fatal("Failed to load .env %q: %v", envPath, err)
		}

		logger.Info("Loaded environment variables from %q", envPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	Conf = &Config{
		GRPCServer: &GRPCServer{
			Host: getEnv("GRPC_HOST", "localhost"),
			Port: getEnvInt("GRPC_PORT", 5001),
		},
		JWT: &JWT{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: getEnvInt("JWT_EXPIRY_SECONDS", 3600),
		},
		GRPCClient: &GRPCClient{
			UserServiceURL: getEnv("GRPC_USER_SERVICE_URL", "localhost:5000"),
			Timeout:        getEnvDuration("GRPC_CLIENT_TIMEOUT_SECONDS", 5),
		},
		Redis: &Redis{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Username: getEnv("REDIS_USERNAME", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Prometheus: &Prometheus{
			METRICS_HOST: getEnv("PROMETHEUS_METRICS_HOST", "localhost"),
			METRICS_PORT: getEnvInt("PROMETHEUS_METRICS_PORT", 5011),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
