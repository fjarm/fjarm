package users

import (
	idempotencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/idempotency/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"errors"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log/slog"
	"testing"
)

func TestUserDomain_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cache := remote.NewFakeRedisCache()
	repo := newInMemoryRepository(logger)
	validator, err := protovalidate.New()
	if err != nil {
		t.Errorf("failed to create a new validator: %v", err)
	}
	dom := newUserDomain(logger, cache, repo, validator)

	tests := map[string]struct {
		reqs []*userspb.CreateUserRequest
		errs []bool
		kind []error
	}{
		"validation_one_valid_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
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
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User:   nil, // The User field is required but is nil here.
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_mismatched_id_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174555")}, // The User ID here ends in 555, which doesn't match the user ID in the field below.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_password_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
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
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
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
					IdempotencyKey: &idempotencypb.IdempotencyKey{},
					UserId:         &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
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
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
		"idempotency_two_identical_id_users": {
			// Two identical users with different idempotency keys should pass without error because the internal state
			// should be hidden from clients.
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key - ends with 888 instead of 999.
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
		"idempotency_two_identical_email_users": {
			// Two identical users with different idempotency keys should pass without error because the internal state
			// should be hidden from clients.
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key - ends with 888 instead of 999.
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")}, // Different user ID from the one in the request above.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")}, // User ID here matches the user ID in the request message.
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
		"idempotency_two_identical_handle_users": {
			// Two identical users with different idempotency keys should pass without error because the internal state
			// should be hidden from clients.
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key - ends with 888 instead of 999.
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")}, // Different user ID from the one in the request above.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")}, // User ID here matches the user ID in the request message.
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false},
			kind: []error{nil, nil},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for i, req := range tc.reqs {
				_, err = dom.createUser(context.Background(), req)
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
		repo.database = map[string]user{}
	}
}
