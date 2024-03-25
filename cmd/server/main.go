package main

import (
	"fmt"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/auth"
	pb "github.com/sagarmaheshwary/microservices-authentication-service/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	if err := env.Load(path.Join(RootDir(), "..", ".env")); err != nil {
		log.Fatalln("Unable to load env!", err)
	}

	host := env.Get("GRPC_HOST", "localhost")
	port, _ := strconv.Atoi(env.Get("GRPC_PORT", "5001"))

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%d", host, port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, auth.NewServer())

	log.Printf("gRPC server started on \"%s:%d\"", host, port)

	grpcServer.Serve(listener)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
