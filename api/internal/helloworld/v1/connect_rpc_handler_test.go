package v1

import (
	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"connectrpc.com/connect"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
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
		expected string
		err      bool
	}{
		"valid_empty_string_input": {
			input:    "",
			expected: "Hello World",
			err:      false,
		},
		"valid_non_empty_string_input": {
			input:    "gleep",
			expected: "Hello World, gleep",
			err:      false,
		},
		"valid_nil_req_and_nil_input": {
			expected: "Hello World",
			err:      false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := &pb.GetHelloWorldRequest{Input: &pb.HelloWorldInput{Input: &tc.input}}

			output, err := client.GetHelloWorld(context.Background(), connect.NewRequest(req))
			if err != nil && !tc.err {
				t.Errorf("GetHelloWorld got an unexpected error: %v", err)
			}

			if output.Msg.GetOutput().GetOutput() != tc.expected {
				t.Errorf("GetHelloWorld got: %v, want: %v", output.Msg.GetOutput().GetOutput(), tc.expected)
			}
		})
	}
}
