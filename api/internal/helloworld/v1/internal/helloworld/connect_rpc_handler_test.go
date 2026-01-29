package helloworld

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/tracing"
)

var srv *httptest.Server = nil

func TestMain(m *testing.M) {
	connectRPCHandler := NewConnectRPCHandler(slog.Default())
	path, handler := helloworldv1connect.NewHelloWorldServiceHandler(connectRPCHandler)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	srv = httptest.NewServer(mux)
	defer srv.Close()

	m.Run()
}

func TestConnectRPCHandler_GetHelloWorld_gRPCClient(t *testing.T) {
	client := helloworldv1connect.NewHelloWorldServiceClient(http.DefaultClient, srv.URL, connect.WithGRPC())

	tests := map[string]struct {
		input    string
		header   []string
		expected string
		err      bool
	}{
		"valid_empty_string_input": {
			input:    "",
			header:   []string{tracing.RequestIDKey, "abc123"},
			expected: "Hello World",
			err:      false,
		},
		"valid_non_empty_string_input": {
			input:    "gleep",
			header:   []string{tracing.RequestIDKey, "abc123"},
			expected: "Hello World, gleep",
			err:      false,
		},
		"valid_nil_req_and_nil_input": {
			header:   []string{tracing.RequestIDKey, "abc123"},
			expected: "Hello World",
			err:      false,
		},
		"invalid_empty_request_id": {
			input:    "gleep",
			header:   []string{tracing.RequestIDKey, ""},
			expected: "",
			err:      true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := &pb.GetHelloWorldRequest{Input: &pb.HelloWorldInput{Input: &tc.input}}
			req := connect.NewRequest(msg)
			req.Header().Set(tc.header[0], tc.header[1])

			output, err := client.GetHelloWorld(context.Background(), req)
			if err != nil && !tc.err {
				t.Errorf("GetHelloWorld got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("GetHelloWorld expected an error but got none: %v", output.Msg.GetOutput().GetOutput())
			}

			if !tc.err && output.Msg.GetOutput().GetOutput() != tc.expected {
				t.Errorf("GetHelloWorld got: %v, want: %v", output.Msg.GetOutput().GetOutput(), tc.expected)
			}
		})
	}
}
