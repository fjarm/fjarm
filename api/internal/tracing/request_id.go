package tracing

import "context"

// RequestIDKey represents the string used as the key to access the request ID value from the request metadata map.
const RequestIDKey = "request-id"

type contextKey string

const requestIDContextKey contextKey = "request-id"

// ContextWithRequestID returns a new context with the request ID value stored.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

// RequestIDFromContext extracts the request ID from the context, returning an empty string if not found.
func RequestIDFromContext(ctx context.Context) string {
	val, _ := ctx.Value(requestIDContextKey).(string)
	return val
}
