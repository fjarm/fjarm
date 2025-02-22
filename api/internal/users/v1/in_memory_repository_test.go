package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
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
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err: []bool{false},
		},
		"valid_two_valid_messages_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
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
					UserId:   &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					Password: &userspb.UserPassword{},
				},
			},
			err: []bool{true},
		},
		"invalid_one_no_email_users": {
			users: []*userspb.User{
				{
					UserId:   &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					Password: &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err: []bool{true},
		},
		"invalid_one_non_uuid_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("user_id")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err: []bool{true},
		},
		"invalid_no_password_no_email_users": {
			users: []*userspb.User{
				{
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
				},
			},
			err: []bool{true},
		},
		"invalid_two_identical_id_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
					Password:     &userspb.UserPassword{Password: proto.String("password1")},
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
