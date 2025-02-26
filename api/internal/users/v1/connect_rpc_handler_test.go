package v1

import (
	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/users/v1/usersv1connect"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"connectrpc.com/connect"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var srv *httptest.Server = nil

func TestMain(m *testing.M) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	connectRPCHandler := NewConnectRPCHandler(logger)
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
		kind []error
	}{}
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
				if !errors.Is(err, tc.kind[index]) {
					t.Errorf("CreateUser got an unexpected error kind: %v", err)
				}
			}
		})
	}
}
