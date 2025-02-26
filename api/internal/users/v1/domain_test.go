package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
	"log/slog"
	"testing"
)

func TestUserDomain_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newInMemoryRepository(logger)
	dom := newUserDomain(logger, repo)

	tests := map[string]struct {
		users []*userspb.User
		errs  []bool
		kind  []error
	}{
		"validation_one_valid_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			errs: []bool{false},
			kind: []error{nil},
		},
		"validation_one_no_password_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{},
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"idempotency_two_distinct_valid_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
		"idempotency_two_identical_id_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for i, msg := range tc.users {
				_, err := dom.createUser(context.Background(), msg)
				if err != nil && !tc.errs[i] {
					t.Errorf("createUser got an unexpected error: %v", err)
				}
				if err == nil && tc.errs[i] {
					t.Errorf("createUser expected an error but got nil at index %d", i)
				}
				if !errors.Is(err, tc.kind[i]) {
					t.Errorf("createUser got an unexpected error type: %v", err)
				}
			}
		})
		// Reset the database for each test run.
		repo.database = map[string]user{}
	}
}
