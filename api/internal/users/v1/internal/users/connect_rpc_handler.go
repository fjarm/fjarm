package users

import (
	"context"
	"log/slog"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
	"buf.build/go/protovalidate"
)

// UserCreator is a sealed interface that is used to create instances of User entities.
type UserCreator interface {
	createUser(ctx context.Context, idemKey string, email string, password string) error
}

// UserCommandManager is an interface that defines mutation operations that can be done to a User entity.
type UserCommandManager interface {
	UserCreator
}

// ConnectRPCHandler defines a ConnectRPC handler for the `fjarm.users.v1.UserService` service.
type ConnectRPCHandler struct {
	usersv1connect.UnimplementedUserServiceHandler
	useCase   UserCreator
	logger    *slog.Logger
	validator protovalidate.Validator
}

// NewConnectRPCHandler creates a concrete users ConnectRPC service with logging and business/domain logic.
func NewConnectRPCHandler(
	logger *slog.Logger,
	uc UserCreator,
	validator protovalidate.Validator,
) *ConnectRPCHandler {
	han := ConnectRPCHandler{
		useCase:   uc,
		logger:    logger,
		validator: validator,
	}
	return &han
}
