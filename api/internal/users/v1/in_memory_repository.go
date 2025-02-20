package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	authentication "github.com/fjarm/fjarm/api/internal/authentication/pkg/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"log/slog"
	"time"
)

const inMemoryRepositoryTag = "in_memory_repository"

type inMemoryRepository struct {
	database  map[string]user
	logger    *slog.Logger
	validator protovalidate.Validator
}

func (repo *inMemoryRepository) createUser(ctx context.Context, msg *userspb.User) (*user, error) {
	logger := repo.logger.With(
		slog.String(logkeys.Tag, inMemoryRepositoryTag),
		slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),
	)
	logger.InfoContext(ctx, "requested user creation")

	// If a user entity with the same ID as the message already exists, return an already exists error.
	_, ok := repo.database[msg.GetUserId().GetUserId()]
	if ok {
		return nil, ErrAlreadyExists
	}

	// The message validation is redundant, but protects against upstream changes in the input/domain layer(s) that
	// should result in invalid input from going uncaught.
	err := repo.validator.Validate(msg)
	if err != nil {
		logger.ErrorContext(ctx,
			"failed to validate user entity",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		// Wrap the error from `protovalidate` so the transport handler can return the correct error code:
		// connect.CodeInvalidArgument.
		return nil, fmt.Errorf("%v: %v", ErrInvalidArgument, err)
	}

	// At this point, the supplied user message should be valid. Convert the Protobuf message to a storage entity.
	// Because `wireUserToStorageUser` returns an error if the message is nil, we don't need to check for `nil` here or
	// elsewhere.
	entity, err := wireUserToStorageUser(msg)
	if err != nil {
		logger.ErrorContext(ctx,
			"failed to convert user message to storage entity",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		// The error message from wireUserToStorageUser is already wrapped with ErrInvalidArgument.
		return &user{}, err
	}

	// This shouldn't happen as the `fjarm.users.v1.UserPassword` message is required when calling
	// `fjarm.users.v1.UserService/CreateUser`. But check for it anyway.
	if !msg.HasPassword() || !msg.GetPassword().HasPassword() {
		logger.ErrorContext(ctx,
			"failed to create user entity with missing credentials",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
		)
		return nil, fmt.Errorf("%v: %v", ErrInvalidArgument, "user message missing password")
	}

	pwd, err := authentication.HashPassword(msg.GetPassword().GetPassword())
	if err != nil {
		logger.ErrorContext(ctx,
			"failed to hash credentials supplied in user message",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		return nil, fmt.Errorf("%v: %v", ErrAuthenticationIssue, err)
	}
	entity.Password = pwd

	// Configure the entity's creation and last updated timestamps.
	now := time.Now()
	entity.CreatedAt = now
	entity.LastUpdated = now

	// Store the entity in the in-memory database.
	repo.database[entity.UserID] = *entity
	return nil, nil
}

func newInMemoryRepository(l *slog.Logger, v protovalidate.Validator) *inMemoryRepository {
	repo := inMemoryRepository{
		database:  map[string]user{},
		logger:    l,
		validator: v,
	}
	return &repo
}
