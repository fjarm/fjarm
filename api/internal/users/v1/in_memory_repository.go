package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"fmt"
	authentication "github.com/fjarm/fjarm/api/internal/authentication/pkg/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"log/slog"
	"time"
)

const inMemoryRepositoryTag = "in_memory_repository"

type inMemoryRepository struct {
	database map[string]user
	logger   *slog.Logger
}

func (repo *inMemoryRepository) createUser(ctx context.Context, msg *userspb.User) (*user, error) {
	logger := repo.logger.With(
		slog.String(logkeys.Tag, inMemoryRepositoryTag),
		slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),
	)
	logger.InfoContext(ctx, "requested user creation")

	if msg == nil {
		return nil, ErrInvalidArgument
	}

	// The message validation is redundant, but protects against upstream changes in the input/domain layer(s) that
	// should result in invalid input from going uncaught.
	err := validateUserMessageForCreate(ctx, msg)
	if err != nil {
		logger.ErrorContext(ctx,
			"failed to validate user message for creation",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		// Wrap the error from `protovalidate` so the transport handler can return the correct error code:
		// connect.CodeInvalidArgument.
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	// If a user entity with the same ID as the message already exists, return an already exists error.
	_, ok := repo.database[msg.GetUserId().GetUserId()]
	if ok {
		return nil, ErrAlreadyExists
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

	pwd, err := authentication.HashPassword(msg.GetPassword().GetPassword())
	if err != nil {
		logger.ErrorContext(ctx,
			"failed to hash credentials supplied in user message",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		return nil, fmt.Errorf("%w: %w", ErrAuthenticationIssue, err)
	}
	entity.Password = pwd

	// Configure the entity's creation and last updated timestamps.
	now := time.Now()
	entity.CreatedAt = now
	entity.LastUpdated = now

	// Store the entity in the in-memory database.
	repo.database[entity.UserID] = *entity
	return &user{}, nil
}

func newInMemoryRepository(l *slog.Logger) *inMemoryRepository {
	repo := inMemoryRepository{
		database: map[string]user{},
		logger:   l,
	}
	return &repo
}
