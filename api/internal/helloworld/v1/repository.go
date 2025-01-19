package v1

import (
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"context"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type repository struct{}

func (repo repository) getHelloWorld(_ context.Context) (*pb.GetHelloWorldResponse, error) {
	msg := "Hello World!"
	res := pb.HelloWorldOutput{
		Output: &msg,
	}
	s := status.Status{
		Code:    int32(codes.OK),
		Message: "OK",
		Details: nil,
	}
	return &pb.GetHelloWorldResponse{
		Status: &s,
		Output: &res,
	}, nil
}
