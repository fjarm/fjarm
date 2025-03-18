package users

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"errors"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing/v1/pkg"
	"log/slog"
	"time"
)

const domainTag = "domain"

type userRepository interface {
	createUser(ctx context.Context, user *userspb.User) (*user, error)
}

// IdempotencyCache is an interface that describes how idempotency is leveraged in the users domain.
type IdempotencyCache interface {
	// Get retrieves the value for a given key if it exists. If the key doesn't exist, a nil slice is returned without
	// any error.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set writes a key value pair in the cache. If the key already exists, an error is returned.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Update writes a key value pair in the cache. If the key already exists, the value is updated. If the key doesn't
	// exist, the key value pair is created.
	Update(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type domain struct {
	logger    *slog.Logger
	repo      userRepository
	validator protovalidate.Validator
}

func newUserDomain(l *slog.Logger, r userRepository, v protovalidate.Validator) userDomain {
	dom := &domain{
		logger:    l,
		repo:      r,
		validator: v,
	}
	return dom
}

// TODO(2025-02-27): Idempotency needs to be handled here after Redis-based caching is introduced.
func (dom *domain) createUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.User, error) {
	logger := dom.logger.With(
		slog.String(logkeys.Tag, domainTag),
		slog.Any(pkg.RequestIDKey, ctx.Value(pkg.RequestIDKey)),
	)

	// Validate the incoming request. The user it contains and its fields will be validated by the repository.
	err := dom.validator.Validate(req)
	// The user ID in the request must match the user ID in the user entity.
	if err != nil || req.GetUserId().GetUserId() != req.GetUser().GetUserId().GetUserId() {
		logger.ErrorContext(ctx,
			"failed to validate incoming request message",
			slog.String(logkeys.Raw, redactedUserMessageString(req.GetUser())),
			slog.Any(logkeys.Err, err),
		)
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	msg := req.GetUser()
	_, err = dom.repo.createUser(ctx, msg)
	if err != nil && errors.Is(err, ErrAlreadyExists) {
		// User creation is idempotent. But, we don't want to leak this information to the client. So, instead of
		// returning the error, we return a successful response without the user's details.
		return &userspb.User{}, nil
	} else if err != nil && errors.Is(err, ErrAuthenticationIssue) {
		// Obscure authentication issues from the client.
		return nil, ErrOperationFailed
	} else if err != nil {
		return nil, err
	}
	// Creating a user is dead simple because enrolling is not the same as authenticating. Users first sign up then
	// log in.
	return &userspb.User{}, nil
}

func (dom *domain) getUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.User, error) {
	return nil, ErrUnimplemented
}

func (dom *domain) updateUser(ctx context.Context, req *userspb.UpdateUserRequest) (*userspb.User, error) {
	return nil, ErrUnimplemented
}

func (dom *domain) deleteUser(ctx context.Context, req *userspb.DeleteUserRequest) error {
	return ErrUnimplemented
}
