package users

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/google/uuid"

	authentication "github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const inMemoryRepositoryTag = "in_memory_repository"

// userStore encapsulates thread-safe in-memory storage for user entities.
type userStore struct {
	mu    sync.RWMutex
	users map[string]user
}

func newUserStore() *userStore {
	return &userStore{
		users: make(map[string]user),
	}
}

// insert atomically checks unique constraints (email, handle, userID) and stores the entity.
func (s *userStore) insert(u user) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.users {
		if existing.EmailAddress == u.EmailAddress {
			return ErrAlreadyExists
		}
		if existing.Handle == u.Handle {
			return ErrAlreadyExists
		}
	}
	if _, ok := s.users[u.UserID]; ok {
		return ErrAlreadyExists
	}

	s.users[u.UserID] = u
	return nil
}

func (s *userStore) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = make(map[string]user)
}

type inMemoryRepository struct {
	store  *userStore
	logger *slog.Logger
}

func (repo *inMemoryRepository) reset() {
	repo.store.reset()
}

func (repo *inMemoryRepository) createUser(ctx context.Context, msg *userspb.User) (*user, error) {
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
	entity.LastUpdated = now

	// Atomically verify unique constraints and store the entity in the in-memory store.
	if err := repo.store.insert(*entity); err != nil {
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
