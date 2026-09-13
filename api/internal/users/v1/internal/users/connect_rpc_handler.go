package users

import (
	"context"
	"log/slog"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
)

type userUseCase interface {
	createUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.User, error)
}

// ConnectRPCHandler defines a ConnectRPC handler for the `fjarm.users.v1.UserService` service.
type ConnectRPCHandler struct {
	usersv1connect.UnimplementedUserServiceHandler
	useCase   userUseCase
	logger    *slog.Logger
	validator protovalidate.Validator
}

// NewConnectRPCHandler creates a concrete users ConnectRPC service with logging and business/domain logic.
func NewConnectRPCHandler(
	logger *slog.Logger,
	domain userUseCase,
	validator protovalidate.Validator,
) *ConnectRPCHandler {
	han := ConnectRPCHandler{
		useCase:   domain,
		logger:    logger,
		validator: validator,
	}
	return &han
}
