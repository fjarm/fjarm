package v1

import (
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"context"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestService_GetHelloWorld(t *testing.T) {
	tests := map[string]struct {
		wantCode   int32
		wantOutput string
	}{
		"valid": {
			wantCode:   int32(codes.OK),
			wantOutput: "Hello World",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			svc := NewService()
			actual, err := svc.GetHelloWorld(context.Background(), &pb.GetHelloWorldRequest{})
			if err != nil {
				t.Errorf("GetHelloWorld got unexpected error: %v", err)
			}
			if tc.wantCode != actual.Status.Code {
				t.Errorf("want %d, got %s", tc.wantCode, actual.Status)
			}
		})
	}
}
