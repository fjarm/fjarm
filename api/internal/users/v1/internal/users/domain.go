package users

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"context"
	"errors"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/cache"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"log/slog"
	"time"
)

const (
	domainTag = "domain"

	// Including the operation (create) in the cache key to differentiate between different idempotent operations. That
	// way, if another operation like update creates the same cache key, it won't collide with the create operation.
	createUserCacheKey = "userservice:create:idempotency:user"
	createUserLockKey  = "userservice:create:lock:user"
	idempotencyKeyTTL  = 24 * time.Hour   // Long-lived to ensure business idempotency
	idempotencyLockTTL = 30 * time.Second // Short-lived to coordinate concurrent requests
)

type userRepository interface {
	createUser(ctx context.Context, user *userspb.User) (*user, error)
}

// idempotencyCache is an interface that describes how idempotency is leveraged in the users domain.
type idempotencyCache interface {
	// Get retrieves the value for a given key if it exists. If the key doesn't exist, a nil slice is returned without
	// any error.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set writes a key value pair in the cache. If the key already exists, an error is returned.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Update writes a key value pair in the cache. If the key already exists, the value is updated. If the key doesn't
	// exist, the key value pair is created.
	Update(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// remoteLock describes methods for acquiring and releasing distributed locks.
type remoteLock interface {
	// AcquireLock creates a lock that can be used to coordinate concurrent requests.
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (string, error)

	// SafeReleaseLock attempts to safely release a lock without interfering with other processes.
	SafeReleaseLock(ctx context.Context, key string, value string) error

	// VerifyLock checks if a process's lock is still valid.
	VerifyLock(cxt context.Context, lockKey string, lockVal string) (bool, error)
}

type domain struct {
	logger    *slog.Logger
	cache     idempotencyCache
	locker    remoteLock
	repo      userRepository
	validator protovalidate.Validator
}

func newUserDomain(
	logger *slog.Logger,
	cache idempotencyCache,
	locker remoteLock,
	repo userRepository,
	validator protovalidate.Validator,
) userDomain {
	dom := &domain{
		logger:    logger,
		cache:     cache,
		locker:    locker,
		repo:      repo,
		validator: validator,
	}
	return dom
}

func (dom *domain) createUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.User, error) {
	logger := dom.logger.With(
		slog.String(logkeys.Tag, domainTag),
		slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),
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

	idempotencyKey := fmt.Sprintf("%s:%s", createUserCacheKey, req.GetIdempotencyKey().GetIdempotencyKey())
	_, err = dom.cache.Get(ctx, idempotencyKey)
	if err == nil {
		// Found a cached response. We can return a successful response without creating the user again.
		return &userspb.User{}, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		// This is a real Redis error, not a cache miss. Return an ErrOperationFailed that will be sent to the client
		// with a connect.CodeInternal status.
		logger.ErrorContext(ctx, "failed to get idempotency key from cache", slog.Any(logkeys.Err, err))
		return nil, ErrOperationFailed
	}

	// At this point, we know that the idempotency key doesn't exist in the cache. So, we can proceed with creating the
	// user. First, we attempt to acquire a distributed lock to indicate to other server replicas that duplicate
	// requests are being processed. In that case, the replicas processing duplicate requests should simply wait for the
	// processing to complete and return the cached result. If the primary server handling the request fails and doesn't
	// update the cache, the replicas should also fail and the client can attempt retrying.
	lockKey := fmt.Sprintf("%s:%s", createUserLockKey, req.GetIdempotencyKey().GetIdempotencyKey())
	lockVal, err := dom.locker.AcquireLock(ctx, lockKey, idempotencyLockTTL)
	// Ensure we release the lock after processing the request.
	defer func() {
		releaseErr := dom.locker.SafeReleaseLock(ctx, lockKey, lockVal)
		if releaseErr != nil {
			logger.ErrorContext(ctx, "failed to release lock in cache", slog.Any(logkeys.Err, releaseErr))
		}
	}()
	if err != nil && !errors.Is(err, cache.ErrKeyExists) {
		// This is a real Redis error, not a cache miss.
		logger.ErrorContext(ctx, "failed to acquire lock in cache", slog.Any(logkeys.Err, err))
		return nil, ErrOperationFailed
	}
	if err != nil && errors.Is(err, cache.ErrKeyExists) {
		// Another process already owns the lock. Retry with a backoff.
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// This is an infinite loop that will only exit when the primary server handling the request caches the result
		// or the timeout is reached.
		for backoff := 50 * time.Millisecond; ; backoff = min(backoff*2, 500*time.Millisecond) {
			select {
			case <-timeoutCtx.Done():
				logger.ErrorContext(ctx, "timed out waiting for concurrent operation to finish")
				return nil, ErrOperationFailed
			default:
				// Try to get the result of the primary server's operation from the cache.
				_, err = dom.cache.Get(ctx, idempotencyKey)
				if err == nil {
					// Found a cached response. We can return a successful response without creating the user again.
					return &userspb.User{}, nil
				} else if !errors.Is(err, cache.ErrCacheMiss) {
					// This is a real Redis error, not a cache miss.
					logger.ErrorContext(ctx, "failed to poll cache for response", slog.Any(logkeys.Err, err))
					return nil, ErrOperationFailed
				}
			}

			select {
			case <-time.After(backoff):
				// Wait before retrying.
			case <-timeoutCtx.Done():
				logger.ErrorContext(ctx, "timed out waiting for concurrent operation to finish")
				return nil, ErrOperationFailed
			}
		}
	}

	// Verify that we have the lock before proceeding.
	verified, err := dom.locker.VerifyLock(ctx, lockKey, lockVal)
	if err != nil || !verified {
		// If there's an error, we do not own the lock (the value in the cache does not match ours), or the lock expired
		// (the key/value pair was deleted), we abort the operation. The client is responsible for retrying. The retry
		// will either find the cached response in the idempotency cache, succeed in creating the user, or find that the
		// user already exists in the database/repository.
		logger.ErrorContext(ctx, "failed to verify lock in cache", slog.Any(logkeys.Err, err))
		return nil, ErrOperationFailed
	}

	msg := req.GetUser()
	_, err = dom.repo.createUser(ctx, msg)
	if err != nil && errors.Is(err, ErrAlreadyExists) {
		logger.WarnContext(ctx, "attempted to create duplicate user", slog.Any(logkeys.Err, err))
		// User creation is idempotent. But, we don't want to leak this information to the client. So, instead of
		// returning the error, we return a successful response without the user's details.
		err = dom.cache.Set(ctx, idempotencyKey, []byte(""), idempotencyKeyTTL)
		if err != nil {
			logger.WarnContext(ctx, "failed to set idempotency key in cache", slog.Any(logkeys.Err, err))
		}
		return &userspb.User{}, nil
	} else if err != nil && errors.Is(err, ErrAuthenticationIssue) {
		// Obscure authentication issues from the client.
		return nil, ErrOperationFailed
	} else if err != nil {
		return nil, err
	}
	// Creating a user is dead simple because enrolling is not the same as authenticating. Users first sign up then
	// log in.
	err = dom.cache.Set(ctx, idempotencyKey, []byte(""), idempotencyKeyTTL)
	if err != nil {
		logger.WarnContext(ctx, "failed to set idempotency key in cache", slog.Any(logkeys.Err, err))
	}
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
