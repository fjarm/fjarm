package v1

import (
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"connectrpc.com/connect"
	"context"
	"log/slog"
	"testing"
)

func TestConnectRPCHandler_GetHelloWorld(t *testing.T) {
	h := NewConnectRPCHandler(slog.Default())
	tests := map[string]struct {
		input    string
		expected string
		err      bool
	}{
		"valid_empty_string_request": {
			input:    "",
			expected: "Hello World",
			err:      false,
		},
		"valid_non_empty_string_request": {
			input:    "gleep glob",
			expected: "Hello World, gleep glob",
			err:      false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := connect.NewRequest(
				&pb.GetHelloWorldRequest{
					Input: &pb.HelloWorldInput{
						Input: &tc.input,
					},
				},
			)
			output, err := h.GetHelloWorld(context.Background(), req)
			if err != nil && !tc.err {
				t.Errorf("GetHelloWorld got an unexpected error: %v", err)
			}
			if output.Msg.GetOutput().GetOutput() != tc.expected {
				t.Errorf("GetHelloWorld got: %v, want: %v", output.Msg.GetOutput().GetOutput(), tc.expected)
			}
		})
	}
}
