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

func TestInMemoryRepository_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newInMemoryRepository(logger)

	tests := map[string]struct {
		users []*userspb.User
		err   []bool
		kind  []error
	}{
		"validation_one_valid_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false},
			kind: []error{nil},
		},
		"validation_one_nil_user": {
			users: []*userspb.User{
				nil,
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_empty_user": {
			users: []*userspb.User{
				{},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_id_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_id_user": {
			users: []*userspb.User{
				{
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_id_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("user_id")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_empty_string_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_contains_spaces_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String(" ")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_email_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_email_user": {
			users: []*userspb.User{
				{
					UserId:   &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName: &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					Handle:   &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password: &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_email_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("gleeper")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_full_name_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("gleeper@glopper.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_full_name_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("gleeper@glopper.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_family_name_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("gleeper@glopper.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_given_name_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{FamilyName: proto.String("foo")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("gleeper@glopper.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_password_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_password_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"idempotency_two_distinct_valid_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
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
			err:  []bool{false, false},
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
			err:  []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
		"idempotency_two_identical_email_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
		"idempotency_two_identical_handle_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
					FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for index, create := range tc.users {
				_, err := repo.createUser(context.Background(), create)
				if err != nil && !tc.err[index] {
					t.Errorf("createUser got an unexpected error: %v", err)
				}
				if err == nil && tc.err[index] {
					t.Errorf("createUser expected an error but got nil")
				}
				if !errors.Is(err, tc.kind[index]) {
					t.Errorf("createUser got an unexpected error type: %v", err)
				}
			}
		})
		// Reset the database for each test run.
		repo.database = map[string]user{}
	}
}
