package main

import (
	"github.com/fjarm/fjarm/api/internal/helloworld"
	pb "github.com/fjarm/fjarm/api/pkg/helloworld/v1"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

const address = "0.0.0.0:8080"

func serve(lis net.Listener) error {
	srv := grpc.NewServer()
	pb.RegisterHelloWorldServiceServer(srv, helloworld.NewService())
	slog.Info("starting server", "addr", address)

	if err := srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}

	if err = serve(lis); err != nil {
		slog.Error("failed to serve", "err", err)
		os.Exit(1)
	}
}
