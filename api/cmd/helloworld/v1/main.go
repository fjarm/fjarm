package main

import (
	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"
	"context"
	"errors"
	"fmt"
	helloworld "github.com/fjarm/fjarm/api/internal/helloworld/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
)

const ip = "[::]"

const mainTag = "main"

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(
		slog.String(logkeys.Service, helloworldv1connect.HelloWorldServiceName),
	)

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		logger.ErrorContext(ctx, "failed to read port from environment", slog.String(logkeys.Tag, mainTag))
		os.Exit(1)
	}
	addr := fmt.Sprintf("%s:%s", ip, port)

	connectRPCHandler := helloworld.NewConnectRPCHandler(logger)
	path, handler := helloworldv1connect.NewHelloWorldServiceHandler(connectRPCHandler)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	srv := &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Handler: mux,
	}

	defer func() {
		logger.InfoContext(ctx, "shut down server", slog.String(logkeys.Tag, mainTag))

		err := srv.Shutdown(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "failed to shut down server", slog.String(logkeys.Tag, mainTag))
		}
	}()

	srvErrChan := make(chan error, 1)
	go func() {
		logger.InfoContext(ctx, "started server", slog.String(logkeys.Tag, mainTag))
		srvErrChan <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err := <-srvErrChan:
		// Error when starting HTTP server.
		logger.ErrorContext(
			ctx,
			"stopping server after error",
			slog.String(logkeys.Tag, mainTag),
			slog.Any(logkeys.Err, err),
		)
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}
}
