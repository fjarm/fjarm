package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"io"
	"log/slog"
	"testing"
)

func TestInMemoryRepository_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&userspb.CreateUserRequest{},
			&userspb.CreateUserResponse{},
			&userspb.GetUserRequest{},
			&userspb.GetUserResponse{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	repo := newInMemoryRepository(logger, validator)

	tests := map[string]struct {
		users []*userspb.User
		err   []bool
	}{
		"valid_empty_slice_users": {
			users: []*userspb.User{},
			err:   []bool{false},
		},
		"valid_one_valid_message_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: ptr("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: ptr("gleeper")},
					Password:     &userspb.UserPassword{Password: ptr("password")},
				},
			},
			err: []bool{false},
		},
		"valid_two_valid_messages_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: ptr("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: ptr("gleeper")},
					Password:     &userspb.UserPassword{Password: ptr("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: ptr("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: ptr("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: ptr("gleeper")},
					Password:     &userspb.UserPassword{Password: ptr("password")},
				},
			},
			err: []bool{false, false},
		},
		"invalid_one_empty_message_users": {
			users: []*userspb.User{
				{},
			},
			err: []bool{true},
		},
		"invalid_one_unset_password_users": {
			users: []*userspb.User{
				{
					UserId:   &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
					Password: &userspb.UserPassword{},
				},
			},
			err: []bool{true},
		},
		"invalid_one_no_email_users": {
			users: []*userspb.User{
				{
					UserId:   &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
					Password: &userspb.UserPassword{Password: ptr("password")},
				},
			},
			err: []bool{true},
		},
		"invalid_one_non_uuid_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: ptr("user_id")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: ptr("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: ptr("gleeper")},
					Password:     &userspb.UserPassword{Password: ptr("password")},
				},
			},
			err: []bool{true},
		},
		"invalid_no_password_no_email_users": {
			users: []*userspb.User{
				{
					UserId: &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
				},
			},
			err: []bool{true},
		},
		"invalid_two_identical_id_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: ptr("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: ptr("gleeper")},
					Password:     &userspb.UserPassword{Password: ptr("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: ptr("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: ptr("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: ptr("gleeper1")},
					Password:     &userspb.UserPassword{Password: ptr("password1")},
				},
			},
			err: []bool{false, true},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for index, create := range tc.users {
				_, err = repo.createUser(context.Background(), create)
				if err != nil && !tc.err[index] {
					t.Errorf("createUser got an unexpected error: %v", err)
				}
				if err == nil && tc.err[index] {
					t.Errorf("createUser expected an error but got nil")
				}
			}
			// Reset the database for each test run.
			repo.database = map[string]user{}
		})
	}
}

func ptr(s string) *string {
	return &s
}
