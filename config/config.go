package config

import (
	"os"
	"strconv"
	"time"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/log"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
	JWT        *jwt
	GRPCClient *grpcClient
	Redis      *Redis
}

type grpcServer struct {
	Host string
	Port int
}

type grpcClient struct {
	UserServiceurl string
	Timeout        time.Duration
}

type jwt struct {
	Secret string
	Expiry int
}

type Redis struct {
	Host     string
	Port     int
	Username string
	Password string
}

func InitConfig() {
	port, err := strconv.Atoi(Getenv("GRPC_PORT", "5001"))

	if err != nil {
		log.Error("Invalid GRPC_PORT value %v", err)
	}

	expiry, err := strconv.Atoi(Getenv("JWT_EXPIRY_SECONDS", "3600"))

	if err != nil {
		log.Error("Invalid JWT_EXPIRY_SECONDS value %v", err)
	}

	timeout, err := strconv.Atoi(Getenv("GRPC_CLIENT_TIMEOUT_SECONDS", "5"))

	if err != nil {
		log.Error("Invalid GRPC_CLIENT_TIMEOUT_SECONDS value %v", err)
	}

	redisPort, err := strconv.Atoi(Getenv("REDIS_PORT", "6379"))

	if err != nil {
		log.Error("Invalid REDIS_PORT value %v", err)
	}

	conf = &Config{
		GRPCServer: &grpcServer{
			Host: Getenv("GRPC_HOST", "localhost"),
			Port: port,
		},
		JWT: &jwt{
			Secret: Getenv("JWT_SECRET", ""),
			Expiry: expiry,
		},
		GRPCClient: &grpcClient{
			UserServiceurl: Getenv("GRPC_USER_SERVICE_URL", "localhost:5000"),
			Timeout:        time.Duration(timeout) * time.Second,
		},
		Redis: &Redis{
			Host:     Getenv("REDIS_HOST", "localhost"),
			Port:     redisPort,
			Username: Getenv("REDIS_USERNAME", ""),
			Password: Getenv("REDIS_PASSWORD", ""),
		},
	}
}

func GetgrpcServer() *grpcServer {
	return conf.GRPCServer
}

func Getjwt() *jwt {
	return conf.JWT
}

func GetgrpcClient() *grpcClient {
	return conf.GRPCClient
}

func GetRedis() *Redis {
	return conf.Redis
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
