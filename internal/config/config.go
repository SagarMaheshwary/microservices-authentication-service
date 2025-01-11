package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
)

var Conf *Config

type Config struct {
	GRPCServer *GRPCServer
	JWT        *JWT
	GRPCClient *GRPCClient
	Redis      *Redis
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

func Init() {
	envPath := path.Join(helper.GetRootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatal("Failed to load .env %q: %v", envPath, err)
	}

	log.Info("Loaded %q", envPath)

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
