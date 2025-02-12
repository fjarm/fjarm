package interceptor

import (
	"context"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

// ErrMetadataNotFound is returned when reading the incoming  request context does not have any metadata.
var ErrMetadataNotFound = status.Error(codes.InvalidArgument, "failed to find metadata")

// ErrRequestIDNotFound is returned when an incoming request does not contain a `request-id` key/value pair.
var ErrRequestIDNotFound = status.Error(codes.InvalidArgument, "failed to find request-id value")

// RequestIDLoggingInterceptor extracts and logs the incoming gRPC request's `request-id` key/value pair in its
// metadata.
func RequestIDLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
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
			slog.String("request-id", reqID),
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
			slog.String("request-id", reqID),
			slog.Duration(logkeys.Duration, duration),
			slog.Any(logkeys.Err, err),
		)

		return res, err
	}
}

func getRequestID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMetadataNotFound
	}

	values := md.Get("request-id")
	if len(values) == 0 {
		return "", ErrRequestIDNotFound
	}

	if values[0] == "" {
		return "", ErrRequestIDNotFound
	}

	return values[0], nil
}
