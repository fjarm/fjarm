package main

import (
	rpc "buf.build/gen/go/fjarm/helloworld/grpc/go/helloworld/v1/helloworldv1grpc"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
)

const ip = "[::]"

func serve(lis net.Listener) error {
	srv := grpc.NewServer()
	rpc.RegisterHelloWorldServiceServer(srv, helloworld.NewService())
	slog.Info("starting server", "addr", lis.Addr())

	reflection.Register(srv)
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
	lis, err := net.Listen("tcp6", addr)
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}
	defer lis.Close()

	if err = serve(lis); err != nil {
		slog.Error("failed to serve", "err", err)
		os.Exit(1)
	}
}
