package interceptor

import (
	"context"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log/slog"
	"time"
)

// RequestIDLoggingGRPCInterceptor extracts and logs the incoming gRPC request's `request-id` key/value pair in its
// metadata. It then verifies that the incoming request does indeed contain a request ID.
//
// If no request ID is found in the incoming request context, the request is immediately rejected.
//
// Note that the provided `logger` must already include a `slog.Handler` that extracts a request ID and adds it as an
// attribute.
func RequestIDLoggingGRPCInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		reqID, err := getRequestID(ctx)

		p, ok := peer.FromContext(ctx)
		clientIP := "unknown"
		if ok {
			clientIP = p.Addr.String()
		}

		l := slog.LevelInfo
		if err != nil {
			l = slog.LevelWarn
		}

		logger.Log(
			ctx,
			l,
			"intercepted request",
			slog.String(tracing.RequestIDKey, reqID),
			slog.Time(logkeys.StartTime, start),
			slog.String(logkeys.Addr, clientIP),
			slog.String(logkeys.Rpc, info.FullMethod),
			slog.Any(logkeys.Err, err),
		)

		var res any
		if err == nil {
			res, err = handler(ctx, req)
		}
		duration := time.Since(start)

		logger.InfoContext(
			ctx,
			"completed request",
			slog.String(tracing.RequestIDKey, reqID),
			slog.Duration(logkeys.Duration, duration),
			slog.Any(logkeys.Err, err),
		)

		return res, err
	}
}

func getRequestID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("failed to find metadata")
	}

	values := md.Get(tracing.RequestIDKey)
	if len(values) == 0 {
		return "", ErrRequestIDNotFound
	}

	if values[0] == "" {
		return "", ErrRequestIDNotFound
	}

	return values[0], nil
}
