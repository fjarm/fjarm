package authenticationv1

import (
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestCreateSessionRequest_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&authenticationpb.CreateSessionRequest{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		request *authenticationpb.CreateSessionRequest
		wantErr bool
	}{
		"valid_create_session_request": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("user@example.com"),
				},
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: false,
		},
		"invalid_empty_idempotency_key": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String(""),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("user@example.com"),
				},
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: true,
		},
		"invalid_non_uuid_idempotency_key": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("invalid-key-format"),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("user@example.com"),
				},
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: true,
		},
		"invalid_missing_idempotency_key": {
			request: &authenticationpb.CreateSessionRequest{
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("user@example.com"),
				},
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: true,
		},
		"invalid_missing_email_address": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: true,
		},
		"invalid_empty_email_address": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String(""),
				},
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: true,
		},
		"invalid_bad_format_email_address": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("notanemail"),
				},
				Password: &userspb.UserPassword{
					Password: proto.String("password123"),
				},
			},
			wantErr: true,
		},
		"invalid_missing_password": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("user@example.com"),
				},
			},
			wantErr: true,
		},
		"invalid_unset_password": {
			request: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress: &userspb.UserEmailAddress{
					EmailAddress: proto.String("user@example.com"),
				},
				Password: &userspb.UserPassword{},
			},
			wantErr: true,
		},
		"invalid_empty_create_session_request": {
			request: &authenticationpb.CreateSessionRequest{},
			wantErr: true,
		},
		"invalid_nil_create_session_request": {
			request: nil,
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.request)
			if err != nil && !tc.wantErr {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
