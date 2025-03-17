package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"connectrpc.com/connect"
	"context"
	"errors"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing/v1/pkg"
	"log/slog"
)

const connectRPCHandlerTag = "connect_rpc_handler"

type userDomain interface {
	createUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.User, error)
	getUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.User, error)
	updateUser(ctx context.Context, req *userspb.UpdateUserRequest) (*userspb.User, error)
	deleteUser(ctx context.Context, req *userspb.DeleteUserRequest) error
}

// ConnectRPCHandler defines a ConnectRPC handler for the `fjarm.users.v1.UserService` service.
type ConnectRPCHandler struct {
	domain    userDomain
	logger    *slog.Logger
	validator protovalidate.Validator
}

// CreateUser handles ConnectRPC requests to create a `User` entity.
func (h *ConnectRPCHandler) CreateUser(
	ctx context.Context,
	req *connect.Request[userspb.CreateUserRequest],
) (*connect.Response[userspb.CreateUserResponse], error) {
	logger := h.logger.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
		slog.Any(pkg.RequestIDKey, ctx.Value(pkg.RequestIDKey)),
	)
	logger.InfoContext(ctx, "received request to create user")

	// Create the user entity.
	usr, err := h.domain.createUser(ctx, req.Msg)
	if err != nil {
		logger.ErrorContext(ctx,
			"failed to create user entity",
			slog.String(logkeys.Raw, req.Msg.String()),
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

	res := &userspb.CreateUserResponse{
		User: usr,
	}
	// User creation was successful.
	return connect.NewResponse(res), nil
}

// GetUser handles ConnectRPC requests to retrieve a `User` entity.
func (h *ConnectRPCHandler) GetUser(
	ctx context.Context,
	req *connect.Request[userspb.GetUserRequest],
) (*connect.Response[userspb.GetUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, ErrUnimplemented)
}

// UpdateUser handles ConnectRPC requests to modify a field in a `User` entity.
func (h *ConnectRPCHandler) UpdateUser(
	ctx context.Context,
	req *connect.Request[userspb.UpdateUserRequest],
) (*connect.Response[userspb.UpdateUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, ErrUnimplemented)
}

// DeleteUser handles ConnectRPC requests to delete an instance of a `User` entity.
func (h *ConnectRPCHandler) DeleteUser(
	ctx context.Context,
	req *connect.Request[userspb.DeleteUserRequest],
) (*connect.Response[userspb.DeleteUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, ErrUnimplemented)
}

// NewConnectRPCHandler creates a concrete users ConnectRPC service with logging and business/domain logic.
func NewConnectRPCHandler(l *slog.Logger) *ConnectRPCHandler {
	logger := l.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
	)

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&userspb.CreateUserRequest{},
			&userspb.CreateUserResponse{},
			&userspb.GetUserRequest{},
			&userspb.GetUserResponse{},
		),
	)

	if err != nil {
		logger.Error(
			"failed to create message validator",
			slog.Any(logkeys.Err, err),
		)
		return nil
	}

	rep := newInMemoryRepository(l)
	dom := newUserDomain(l, rep, validator)
	han := ConnectRPCHandler{
		domain:    dom,
		logger:    l,
		validator: validator,
	}
	return &han
}
