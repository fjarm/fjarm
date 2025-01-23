package main

import (
	rpc "buf.build/gen/go/fjarm/fjarm/grpc/go/fjarm/helloworld/v1/helloworldv1grpc"
	"context"
	"fmt"
	helloworld "github.com/fjarm/fjarm/api/internal/helloworld/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const ip = "[::]"

const serviceName = "helloworld"

func main() {
	// Use flags here instead of environment variable
	port := os.Getenv("PORT")
	if port == "" {
		slog.Error(
			"failed to read port from environment",
			slog.String(logkeys.Service, serviceName),
		)
		os.Exit(1)
	}

	ctx := context.Background()

	addr := fmt.Sprintf("%s:%s", ip, port)
	lis, err := net.Listen("tcp6", addr)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to create listener",
			slog.String(logkeys.Service, serviceName),
			slog.String(logkeys.Addr, addr),
			slog.String(logkeys.Err, err.Error()),
		)
		os.Exit(1)
	}

	srv := grpc.NewServer()
	han := helloworld.NewGrpcHandler()
	if han == nil {
		slog.ErrorContext(ctx, "failed to create gRPC handler", slog.String(logkeys.Service, serviceName))
		os.Exit(1)
	}
	rpc.RegisterHelloWorldServiceServer(srv, han)
	reflection.Register(srv)

	closer := func() {
		slog.InfoContext(ctx, "stopping server", slog.String(logkeys.Service, serviceName))
		e := lis.Close()
		if e != nil {
			slog.ErrorContext(ctx, "failed to close listener", slog.String(logkeys.Service, serviceName))
		}
		srv.GracefulStop()
	}
	defer closer()

	slog.InfoContext(
		ctx,
		"starting server",
		slog.String(logkeys.Service, serviceName),
		slog.String(logkeys.Addr, addr),
	)

	go func() {
		e := srv.Serve(lis)
		if e != nil {
			slog.ErrorContext(
				ctx,
				"failed to start serving",
				slog.String(logkeys.Service, serviceName),
				slog.String(logkeys.Err, err.Error()),
			)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
