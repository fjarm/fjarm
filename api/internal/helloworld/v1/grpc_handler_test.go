package v1

import (
	rpc "buf.build/gen/go/fjarm/fjarm/grpc/go/fjarm/helloworld/v1/helloworldv1grpc"
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const (
	bufSize = 1024 * 1024
)

func TestGrpcHandler_GetHelloWorld(t *testing.T) {
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()

	closer := func() {
		err := listener.Close()
		if err != nil {
			t.Logf("failed to close buf conn: %v", err)
		}
		server.GracefulStop()
	}
	defer closer()

	handler := NewGrpcHandler()
	rpc.RegisterHelloWorldServiceServer(server, handler)

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		t.Fatalf("failed to create new client with error: %v", err)
	}

	client := rpc.NewHelloWorldServiceClient(conn)

	tests := map[string]struct {
		given string
		want  string
		err   bool
	}{
		"valid_empty_request": {
			given: "",
			want:  "Hello World",
			err:   false,
		},
		"valid_non_empty_request": {
			given: "gleep",
			want:  "Hello World, gleep",
			err:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			req := &pb.GetHelloWorldRequest{Input: &pb.HelloWorldInput{Input: &tc.given}}

			resp, e := client.GetHelloWorld(ctx, req)
			if e != nil && !tc.err {
				t.Errorf("GetHelloWorld got an unexpected error: %v", e)
			}

			if resp.GetOutput().GetOutput() != tc.want {
				t.Errorf("GetHelloWorld got: %v, want: %v", resp, tc.want)
			}
		})
	}
}
