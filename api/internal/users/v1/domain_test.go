package v1

import (
	idempotencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/idempotency/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"errors"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log/slog"
	"testing"
)

func TestUserDomain_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newInMemoryRepository(logger)
	validator, err := protovalidate.New()
	if err != nil {
		t.Errorf("failed to create a new validator: %v", err)
	}
	dom := newUserDomain(logger, repo, validator)

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
		"validation_one_no_password_user": {
			reqs: []*userspb.CreateUserRequest{
				{
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
		"idempotency_two_distinct_valid_users": {
			reqs: []*userspb.CreateUserRequest{
				{
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
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
			reqs: []*userspb.CreateUserRequest{
				{
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
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
