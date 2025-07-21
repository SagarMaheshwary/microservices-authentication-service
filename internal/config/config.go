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
	Jaeger     *Jaeger
}

type GRPCServer struct {
	Host string
	Port int
}

type GRPCClient struct {
	UserServiceURL string
	TimeoutSeconds time.Duration
}

type JWT struct {
	Secret        string
	ExpirySeconds time.Duration
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
			Secret:        getEnv("JWT_SECRET", ""),
			ExpirySeconds: getEnvDurationSeconds("JWT_EXPIRY_SECONDS", 3600),
		},
		GRPCClient: &GRPCClient{
			UserServiceURL: getEnv("GRPC_USER_SERVICE_URL", "localhost:5000"),
			TimeoutSeconds: getEnvDurationSeconds("GRPC_CLIENT_TIMEOUT_SECONDS", 5),
		},
		Redis: &Redis{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Username: getEnv("REDIS_USERNAME", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Prometheus: &Prometheus{
			URL: getEnv("PROMETHEUS_URL", "localhost:5011"),
		},
		Jaeger: &Jaeger{
			URL: getEnv("JAEGER_URL", "localhost:4318"),
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

func getEnvDurationSeconds(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
