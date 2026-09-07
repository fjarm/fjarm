package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/authentication/v1/authenticationv1connect"
	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/internal/authentication"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	obfuscation "github.com/fjarm/fjarm/api/internal/obfuscation/v1/pkg/interceptor"
	tracing "github.com/fjarm/fjarm/api/internal/tracing/v1/pkg/interceptor"
)

const ip = "[::]"

const mainTag = "main"

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(
		slog.String(logkeys.Service, authenticationv1connect.AuthenticationServiceName),
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

	interceptors := connect.WithInterceptors(
		obfuscation.NewConnectRPCConstantTimingInterceptor(logger, obfuscation.DelayDuration_100ms),
		tracing.NewConnectRPCRequestIDLoggingInterceptor(logger),
	)
	connectRPCHandler, err := authentication.NewConnectRPCHandler(logger)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to initialize ConnectRPC handler",
			slog.String(logkeys.Tag, mainTag),
			slog.Any(logkeys.Err, err),
		)
	}
	path, handler := authenticationv1connect.NewAuthenticationServiceHandler(connectRPCHandler, interceptors)

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

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "failed to shut down server", slog.String(logkeys.Tag, mainTag), slog.Any(logkeys.Err, err))
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
