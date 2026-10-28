package users

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
	"github.com/fjarm/fjarm/api/internal/cache"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const (
	createUserUseCaseTag = "create_user_use_case"

	// Including the operation (create) in the cache key to differentiate between different idempotent operations. That
	// way, if another operation like update creates the same cache key, it won't collide with the create operation.
	createUserCacheKey = "userservice:create:idempotency:user"
	createUserLockKey  = "userservice:create:lock:user"
	idempotencyKeyTTL  = 24 * time.Hour   // Long-lived to ensure business idempotency
	idempotencyLockTTL = 30 * time.Second // Short-lived to coordinate concurrent requests
)

// UserRepository orchestrates persisting a User entity to durable storage like a SQL database.
type UserRepository interface {
	create(ctx context.Context, user *User) error
}

// IdempotencyKeyStore is a data store that contains idempotency keys. The interface can be used to verify if an
// idempotency key already exists - i.e. a request has already been served/handled. The interface can also be used to
// create an entry in the underlying data store - i.e. cache a newly handled request's response.
type IdempotencyKeyStore interface {
	isDuplicate(ctx context.Context, key string) (bool, error)
	saveResult(ctx context.Context, key string) error
}

// DistributedLocker handles distributed mutual exclusion for user resources. It can be used to verify that an instance
// of a users is the sole instance responsible for a particular request.
type DistributedLocker interface {
	lock(ctx context.Context, key string) (string, error)
	unlock(ctx context.Context, key string, value string) error
	verifyLockAcquired(cxt context.Context, lockKey string, lockVal string) (bool, error)
}

// UserUseCase orchestrates business logic for user-related operations like creating a user, reading a user, updating a
// user, and deleting a user.
type UserUseCase struct {
	logger       *slog.Logger
	idemKeyStore IdempotencyKeyStore
	distLock     DistributedLocker
	repo         UserRepository
}

// NewUserUseCase creates an instance of UserUseCase.
func NewUserUseCase(
	logger *slog.Logger,
	idemKeyStore IdempotencyKeyStore,
	distLock DistributedLocker,
	repo UserRepository,
) *UserUseCase {
	return &UserUseCase{
		logger:       logger,
		idemKeyStore: idemKeyStore,
		distLock:     distLock,
		repo:         repo,
	}
}

func (uc *UserUseCase) createUser(ctx context.Context, idemKey string, email string, password string) error {
	logger := uc.logger.With(
		slog.String(logkeys.Tag, createUserUseCaseTag),
		slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
	)

	// Validate inputs before cache access and distributed locking
	if strings.TrimSpace(idemKey) == "" {
		return fmt.Errorf("%w: idempotency key is required", ErrInvalidArgument)
	}
	if _, err := uuid.Parse(idemKey); err != nil {
		return fmt.Errorf("%w: idempotency key must be a valid UUID", ErrInvalidArgument)
	}
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("%w: email address is required", ErrInvalidArgument)
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%w: invalid email address", ErrInvalidArgument)
	}
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("%w: password is required", ErrInvalidArgument)
	}

	idempotencyKey := fmt.Sprintf("%s:%s", createUserCacheKey, idemKey)
	dup, err := uc.idemKeyStore.isDuplicate(ctx, idempotencyKey)
	if err != nil {
		return err
	}
	if dup {
		// Found a cached response. We can return a successful response without creating the user again.
		return nil
	}

	// At this point, we know that the idempotency key doesn't exist in the cache. So, we can proceed with creating the
	// user. First, we attempt to acquire a distributed lock to indicate to other server replicas that duplicate
	// requests are being processed. In that case, the replicas processing duplicate requests should simply wait for the
	// processing to complete and return the cached result. If the primary server handling the request fails and doesn't
	// update the cache, the replicas should also fail and the client can attempt retrying.
	lockKey := fmt.Sprintf("%s:%s", createUserLockKey, idemKey)
	lockVal, err := uc.distLock.lock(ctx, lockKey)
	// Ensure we release the lock after processing the request.
	defer func() {
		if lockVal == "" {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		releaseErr := uc.distLock.unlock(cleanupCtx, lockKey, lockVal)
		if releaseErr != nil {
			logger.ErrorContext(ctx, "failed to release lock in cache", slog.Any(logkeys.Err, releaseErr))
		}
	}()

	if err != nil && !errors.Is(err, cache.ErrKeyExists) {
		// This is a real Redis error, not a cache miss.
		logger.ErrorContext(ctx, "failed to acquire lock in cache", slog.Any(logkeys.Err, err))
		return ErrOperationFailed
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
				return ErrOperationFailed
			default:
				// Try to get the result of the primary server's operation from the cache.
				dup, err = uc.idemKeyStore.isDuplicate(ctx, idempotencyKey)
				if err != nil {
					return err
				}
				if dup {
					// Found a cached response. We can return a successful response without creating the user again.
					return nil
				}
			}

			select {
			case <-time.After(backoff):
				// Wait before retrying.
			case <-timeoutCtx.Done():
				logger.ErrorContext(ctx, "timed out waiting for concurrent operation to finish")
				return ErrOperationFailed
			}
		}
	}

	// Verify that we have the lock before proceeding.
	verified, err := uc.distLock.verifyLockAcquired(ctx, lockKey, lockVal)
	if err != nil || !verified {
		// If there's an error, we do not own the lock (the value in the cache does not match ours), or the lock expired
		// (the key/value pair was deleted), we abort the operation. The client is responsible for retrying. The retry
		// will either find the cached response in the idempotency cache, succeed in creating the user, or find that the
		// user already exists in the database/repository.
		logger.ErrorContext(ctx, "failed to verify lock in cache", slog.Any(logkeys.Err, err))
		return ErrOperationFailed
	}

	parts := strings.Split(addr.Address, "@")
	localName := parts[0]
	h := fnv.New64a()
	h.Write([]byte(addr.Address))
	username := fmt.Sprintf("%s%d", localName, h.Sum64())

	pwd, err := passwords.HashPassword(password)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to hash credentials supplied in user message",
			slog.Any(logkeys.Err, err),
		)
		return ErrOperationFailed
	}

	user := User{
		UserID:       uuid.NewString(),
		EmailAddress: addr.Address,
		Password:     pwd,
		Handle:       username,
	}

	err = uc.repo.create(ctx, &user)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create user in repository", slog.Any(logkeys.Err, err))
		return err
	}
	// Creating a user is dead simple because enrolling is not the same as authenticating. Users first sign up then
	// log in.
	err = uc.idemKeyStore.saveResult(ctx, idempotencyKey)
	if err != nil {
		logger.WarnContext(ctx, "failed to set idempotency key in cache", slog.Any(logkeys.Err, err))
	}
	return nil
}
