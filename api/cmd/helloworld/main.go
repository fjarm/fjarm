package main

import (
	"fmt"
	"github.com/fjarm/fjarm/api/internal/helloworld"
	pb "github.com/fjarm/fjarm/api/pkg/helloworld/v1"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

const ip = "0.0.0.0"

func serve(lis net.Listener) error {
	srv := grpc.NewServer()
	pb.RegisterHelloWorldServiceServer(srv, helloworld.NewService())
	slog.Info("starting server", "addr", lis.Addr())

	if err := srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("failed to read port from environment")
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%s", ip, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}

	if err = serve(lis); err != nil {
		slog.Error("failed to serve", "err", err)
		os.Exit(1)
	}
}
