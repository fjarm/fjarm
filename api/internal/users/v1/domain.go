package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
)

type domain struct {
}

func newUserDomain() userDomain {
	dom := &domain{}
	return dom
}

func (dom *domain) createUser(ctx context.Context, user *userspb.User) (*userspb.User, error) {
	return nil, ErrUnimplemented
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
