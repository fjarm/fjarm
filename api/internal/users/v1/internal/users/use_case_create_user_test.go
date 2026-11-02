package users

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"google.golang.org/protobuf/proto"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

func TestUserUseCase_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cache := remote.NewFakeRedisCache()
	repo := newUserStore()
	idemStore := newFakeRedisIdempotencyKeyStore(cache)
	distLocker := newFakeRedisDistributedLocker(cache)
	uc := NewUserUseCase(logger, idemStore, distLocker, repo)

	tests := map[string]struct {
		reqs []*userspb.CreateUserRequest
		errs []bool
		kind []error
	}{
		"validation_one_valid_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false},
			kind: []error{nil},
		},
		"validation_one_nil_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User:           nil, // The User field is required but is nil here.
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_password_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{},
					},
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_idempotency_key_request": {
			reqs: []*userspb.CreateUserRequest{
				{
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_idempotency_key_request": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String(""),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"idempotency_two_distinct_valid_users": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper2")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
		"idempotency_two_identical_email_users": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key - ends with 888 instead of 999.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")}, // User ID here matches the user ID in the request message.
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
		"idempotency_retry_identical_email_user_with_same_idempotency_key": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key triggers ErrAlreadyExists and caches the key.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Retrying with the same idempotency key hits the idempotency cache.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, true, false},
			kind: []error{nil, ErrAlreadyExists, nil},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for i, req := range tc.reqs {
				err := uc.createUser(context.Background(), req.GetIdempotencyKey(), req.GetUser().GetEmailAddress().GetEmailAddress(), req.GetUser().GetPassword().GetPassword())
				if err != nil && !tc.errs[i] {
					t.Errorf("createUser got an unexpected error: %v", err)
				}
				if err == nil && tc.errs[i] {
					t.Errorf("createUser expected an error but got nil")
				}
				if !errors.Is(err, tc.kind[i]) {
					t.Errorf("createUser got an unexpected error type: %v", err)
				}
			}
		})
		// Reset the database for each test run.
		cache.Reset()
		repo.reset()
	}
}
