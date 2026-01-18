package users

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/users/v1/usersv1connect"
	idempotencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/idempotency/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
	"github.com/fjarm/fjarm/api/internal/logkeys"
)

var srv *httptest.Server = nil

func TestMain(m *testing.M) {
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
		logger.Error("failed to create message validator", slog.Any(logkeys.Err, err))
		os.Exit(1)
	}

	cache := remote.NewFakeRedisCache()
	repo := newInMemoryRepository(logger)
	dom := newUserDomain(logger, cache, cache, repo, validator)
	connectRPCHandler := NewConnectRPCHandler(logger, dom, validator)
	path, handler := usersv1connect.NewUserServiceHandler(connectRPCHandler)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	srv = httptest.NewServer(mux)
	defer srv.Close()

	os.Exit(m.Run())
}

func TestConnectRPCHandler_CreateUser_gRPCClient(t *testing.T) {
	tests := map[string]struct {
		reqs []*userspb.CreateUserRequest
		errs []bool
		code []connect.Code
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
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false},
			code: []connect.Code{},
		},
		"validation_one_nil_user": {
			reqs: []*userspb.CreateUserRequest{
				nil,
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
		},
		"validation_one_no_idempotency_key_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{},
					UserId:         &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
		},
		"validation_one_no_id_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: &idempotencypb.IdempotencyKey{
						IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
						Timestamp:      timestamppb.Now(),
					},
					UserId: &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					User: &userspb.User{
						UserId:       &userspb.UserId{},
						FullName:     &userspb.UserFullName{GivenName: proto.String("foo"), FamilyName: proto.String("bar")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
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
			code: []connect.Code{},
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
			code: []connect.Code{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := usersv1connect.NewUserServiceClient(http.DefaultClient, srv.URL, connect.WithGRPC())

			for index, req := range tc.reqs {
				_, err := client.CreateUser(context.Background(), connect.NewRequest(req))
				if err != nil && !tc.errs[index] {
					t.Errorf("CreateUser got an unexpected error: %v", err)
				}
				if err == nil && tc.errs[index] {
					t.Errorf("CreateUser expected an error but got nil")
				}
				if tc.errs[index] && connect.CodeOf(err) != tc.code[index] {
					t.Errorf("CreateUser got an unexpected error kind: %v, wanted: %v", err, tc.code[index])
				}
			}
		})
	}
}
