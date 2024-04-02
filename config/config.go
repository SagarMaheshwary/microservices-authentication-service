package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
	JWT        *jwt
	GRPCClient *grpcClient
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

func InitConfig() {
	port, err := strconv.Atoi(Getenv("GRPC_PORT", "5001"))

	if err != nil {
		log.Println("Invalid GRPC_PORT value", err)
	}

	expiry, err := strconv.Atoi(Getenv("JWT_EXPIRY_SECONDS", "3600"))

	if err != nil {
		log.Println("Invalid JWT_EXPIRY_SECONDS value", err)
	}

	timeout, err := strconv.Atoi(Getenv("GRPC_CLIENT_TIMEOUT_SECONDS", "5"))

	if err != nil {
		log.Println("Invalid GRPC_CLIENT_TIMEOUT_SECONDS value", err)
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

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
