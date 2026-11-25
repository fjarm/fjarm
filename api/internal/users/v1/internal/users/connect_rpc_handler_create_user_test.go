package users

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
	"github.com/fjarm/fjarm/api/internal/logkeys"
)

var (
	srv       *httptest.Server = nil
	buf       bytes.Buffer
	pwdHasher *fakePasswordHasher
)

func TestMain(m *testing.M) {
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&userspb.CreateUserRequest{},
			&userspb.CreateUserResponse{},
		),
	)
	if err != nil {
		logger.Error("failed to create message validator", slog.Any(logkeys.Err, err))
		os.Exit(1)
	}

	cache := remote.NewFakeRedisCache()
	repo := newUserStore()
	idemStore := newFakeRedisIdempotencyKeyStore(cache)
	distLocker := newFakeRedisDistributedLocker(cache)
	pwdHasher = newFakePasswordHasher()
	dom := NewUserUseCase(logger, idemStore, distLocker, repo, pwdHasher)
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
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
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
		"validation_one_nil_user_field": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User:           nil,
				},
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
		},
		"validation_one_unset_idempotency_key_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
		},
		"validation_one_no_idempotency_key_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String(""),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
		},
		"validation_one_no_password_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{},
					},
				},
			},
			errs: []bool{true},
			code: []connect.Code{connect.CodeInvalidArgument},
		},
		"validation_one_no_id_user": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{},
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
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key - ends with 888 instead of 999.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo2@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper2")},
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
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174999"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"), // Different idempotency key - ends with 888 instead of 999.
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false},
			code: []connect.Code{},
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
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174888"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
						Password:     &userspb.UserPassword{Password: proto.String("password")},
					},
				},
			},
			errs: []bool{false, false, false},
			code: []connect.Code{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := usersv1connect.NewUserServiceClient(http.DefaultClient, srv.URL, connect.WithGRPC())

			for index, req := range tc.reqs {
				_, err := client.CreateUser(context.Background(), req)
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

func TestConnectRPCHandler_CreateUser_LogOutput(t *testing.T) {
	tests := map[string]struct {
		reqs     []*userspb.CreateUserRequest
		expected []string
		excluded []string
		errs     []bool
	}{
		"valid_one_user_request": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174701"),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174001")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("log1@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper_log1")},
						Password:     &userspb.UserPassword{Password: proto.String("secret_password_1")},
					},
				},
			},
			expected: []string{"INFO", "msg=\"received request to create user\""},
			excluded: []string{"secret_password_1", "log1@bar.com", "gleeper_log1"},
			errs:     []bool{false},
		},
		"invalid_empty_idempotency_key_request": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String(""),
					User: &userspb.User{
						UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174002")},
						EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("log2@bar.com")},
						Handle:       &userspb.UserHandle{Handle: proto.String("gleeper_log2")},
						Password:     &userspb.UserPassword{Password: proto.String("secret_password_2")},
					},
				},
			},
			expected: []string{
				"ERROR",
				"msg=\"failed to validate incoming request message\"",
				"123e4567-e89b-12d3-a456-426614174002",
			},
			excluded: []string{"secret_password_2", "log2@bar.com", "gleeper_log2"},
			errs:     []bool{true},
		},
		"invalid_missing_email_address_request": {
			reqs: []*userspb.CreateUserRequest{
				{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174704"),
					User: &userspb.User{
						UserId:   &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174004")},
						Handle:   &userspb.UserHandle{Handle: proto.String("gleeper_log4")},
						Password: &userspb.UserPassword{Password: proto.String("secret_password_4")},
					},
				},
			},
			expected: []string{
				"ERROR",
				"msg=\"failed to create user entity\"",
				"123e4567-e89b-12d3-a456-426614174004",
				"error=\"supplied argument is invalid: email address is required\"",
			},
			excluded: []string{"secret_password_4", "gleeper_log4"},
			errs:     []bool{true},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			buf.Reset()
			client := usersv1connect.NewUserServiceClient(http.DefaultClient, srv.URL, connect.WithGRPC())

			for index, req := range tc.reqs {
				_, err := client.CreateUser(context.Background(), req)
				if err != nil && !tc.errs[index] {
					t.Errorf("CreateUser got an unexpected error: %v", err)
				}
				if err == nil && tc.errs[index] {
					t.Errorf("CreateUser expected an error but got nil")
				}
			}

			actual := buf.String()
			for _, exp := range tc.expected {
				if !strings.Contains(actual, exp) {
					t.Errorf("CreateUser log output got: %v, want: %v", actual, tc.expected)
				}
			}
			for _, exc := range tc.excluded {
				if strings.Contains(actual, exc) {
					t.Errorf("CreateUser log output got: %v, must exclude: %v", actual, tc.excluded)
				}
			}
		})
	}
}

func TestConnectRPCHandler_CreateUser_HashingFailure(t *testing.T) {
	tests := map[string]struct {
		givenReq          *userspb.CreateUserRequest
		givenHashErr      error
		wantErr           bool
		wantCode          connect.Code
		wantErrExcludes   string
		wantLogContains   []string
		wantLogExcludes   []string
	}{
		"hashing_failure_returns_internal_error_and_redacts_logs": {
			givenReq: &userspb.CreateUserRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174705"),
				User: &userspb.User{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174005")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("hashfail@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper_hashfail")},
					Password:     &userspb.UserPassword{Password: proto.String("secret_password_5")},
				},
			},
			givenHashErr:    errors.New("simulated hashing failure"),
			wantErr:         true,
			wantCode:        connect.CodeInternal,
			wantErrExcludes: "simulated hashing failure",
			wantLogContains: []string{
				"ERROR",
				"msg=\"failed to hash credentials supplied in user message\"",
				"msg=\"failed to create user entity\"",
			},
			wantLogExcludes: []string{"secret_password_5", "hashfail@bar.com", "gleeper_hashfail"},
		},
		"valid_hashing_returns_ok_and_redacts_logs": {
			givenReq: &userspb.CreateUserRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174706"),
				User: &userspb.User{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174006")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("hashok@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper_hashok")},
					Password:     &userspb.UserPassword{Password: proto.String("secret_password_6")},
				},
			},
			givenHashErr:    nil,
			wantErr:         false,
			wantLogContains: []string{"INFO", "msg=\"received request to create user\""},
			wantLogExcludes: []string{"secret_password_6", "hashok@bar.com", "gleeper_hashok"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			buf.Reset()
			pwdHasher.setErr(tc.givenHashErr)
			t.Cleanup(pwdHasher.reset)
			client := usersv1connect.NewUserServiceClient(http.DefaultClient, srv.URL, connect.WithGRPC())

			// When
			_, err := client.CreateUser(context.Background(), tc.givenReq)

			// Then
			if tc.wantErr && err == nil {
				t.Fatal("CreateUser expected an error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("CreateUser got an unexpected error: %v", err)
			}
			if tc.wantErr && connect.CodeOf(err) != tc.wantCode {
				t.Errorf("CreateUser got code %v, want %v", connect.CodeOf(err), tc.wantCode)
			}
			if tc.wantErrExcludes != "" && err != nil && strings.Contains(err.Error(), tc.wantErrExcludes) {
				t.Errorf("CreateUser leaked internal hashing error %q to client: %v", tc.wantErrExcludes, err)
			}

			logOutput := buf.String()
			for _, expected := range tc.wantLogContains {
				if !strings.Contains(logOutput, expected) {
					t.Errorf("CreateUser log output missing %q; got: %s", expected, logOutput)
				}
			}
			for _, excluded := range tc.wantLogExcludes {
				if strings.Contains(logOutput, excluded) {
					t.Errorf("CreateUser log output leaked sensitive field %q; got: %s", excluded, logOutput)
				}
			}
		})
	}
}
