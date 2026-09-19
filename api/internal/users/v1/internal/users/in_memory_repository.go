package users

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/google/uuid"

	authentication "github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const inMemoryRepositoryTag = "in_memory_repository"

type inMemoryRepository struct {
	store  *userStore
	logger *slog.Logger
}

func (repo *inMemoryRepository) reset() {
	repo.store.reset()
}

func (repo *inMemoryRepository) createUser(ctx context.Context, msg *userspb.User) (*User, error) {
	logger := repo.logger.With(
		slog.String(logkeys.Tag, inMemoryRepositoryTag),
		slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
	)
	logger.InfoContext(ctx, "requested user creation")

	if msg == nil {
		return nil, ErrInvalidArgument
	}

	// The message validation is redundant, but protects against upstream changes in the input/domain layer(s) that
	// should result in invalid input from going uncaught.
	err := validateUserMessageForCreate(msg)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to validate user message for creation",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		// Wrap the error from `protovalidate` so the transport handler can return the correct error code:
		// connect.CodeInvalidArgument.
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	// At this point, the supplied user message should be valid. Convert the Protobuf message to a storage entity.
	// Because `wireUserToStorageUser` returns an error if the message is nil, we don't need to check for `nil` here or
	// elsewhere.
	entity, err := wireUserToStorageUser(msg)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to convert user message to storage entity",
			slog.String(logkeys.Raw, redactedUserMessageString(msg)),
			slog.Any(logkeys.Err, err),
		)
		// The error message from wireUserToStorageUser is already wrapped with ErrInvalidArgument.
		return nil, err
	}
	entity.UserID = uuid.NewString()

	pwd, err := authentication.HashPassword(msg.GetPassword().GetPassword())
	if err != nil {
		logger.ErrorContext(
			ctx,
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
	entity.UpdatedAt = now

	// Atomically verify unique constraints and store the entity in the in-memory store.
	if err := repo.store.create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func newInMemoryRepository(l *slog.Logger) *inMemoryRepository {
	repo := inMemoryRepository{
		store:  newUserStore(),
		logger: l,
	}
	return &repo
}
