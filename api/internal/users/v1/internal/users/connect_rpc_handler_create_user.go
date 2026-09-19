package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"

	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const connectRPCHandlerCreateUserTag = "connect_rpc_handler_create_user"

// CreateUser handles ConnectRPC requests to create a `User` entity.
func (h *ConnectRPCHandler) CreateUser(
	ctx context.Context,
	req *userspb.CreateUserRequest,
) (*userspb.CreateUserResponse, error) {
	logger := h.logger.With(
		slog.String(logkeys.Tag, connectRPCHandlerCreateUserTag),
		slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
	)
	logger.InfoContext(ctx, "received request to create user")

	if req == nil {
		logger.ErrorContext(ctx, "request cannot be nil")
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: request cannot be nil", ErrInvalidArgument))
	}

	err := h.validator.Validate(req)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to validate incoming request message",
			slog.String(logkeys.Raw, redactedUserMessageString(req.GetUser())),
			slog.Any(logkeys.Err, err),
		)
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: request cannot be invalid", ErrInvalidArgument))
	}

	// Create the user entity.
	err = h.useCase.CreateUser(ctx, req.GetIdempotencyKey(), req.GetUser().GetEmailAddress().GetEmailAddress(), req.GetUser().GetPassword().GetPassword())
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to create user entity",
			slog.String(logkeys.Raw, req.String()),
			slog.Any(logkeys.Err, err),
		)
	}
	if err != nil && errors.Is(err, ErrOperationFailed) {
		// Typically returned for authentication issues. Obscure this from the client.
		return nil, connect.NewError(connect.CodeInternal, err)
	} else if err != nil && errors.Is(err, ErrInvalidArgument) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		// A true edge case. This shouldn't happen.
		return nil, connect.NewError(connect.CodeUnknown, ErrOperationFailed)
	}

	// User creation was successful.
	return &userspb.CreateUserResponse{}, nil
}
