package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"errors"
	"log/slog"
)

type userRepository interface {
	createUser(ctx context.Context, user *userspb.User) (*user, error)
}

type domain struct {
	logger *slog.Logger
	repo   userRepository
}

func newUserDomain(l *slog.Logger, r userRepository) userDomain {
	dom := &domain{
		logger: l,
		repo:   r,
	}
	return dom
}

func (dom *domain) createUser(ctx context.Context, user *userspb.User) (*userspb.User, error) {
	_, err := dom.repo.createUser(ctx, user)
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

func (dom *domain) getUserWithID(ctx context.Context, id *userspb.UserId) (*userspb.User, error) {
	return nil, ErrUnimplemented
}

func (dom *domain) updateUser(ctx context.Context, user *userspb.User) (*userspb.User, error) {
	return nil, ErrUnimplemented
}

func (dom *domain) deleteUser(ctx context.Context, user *userspb.User) error {
	return ErrUnimplemented
}
