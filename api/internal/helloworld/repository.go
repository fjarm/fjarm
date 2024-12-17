package helloworld

import (
	pb "buf.build/gen/go/fjarm/helloworld/protocolbuffers/go/helloworld/v1"
	"context"
	"google.golang.org/grpc/codes"
)

type repository struct{}

func (repo repository) getHelloWorld(_ context.Context) (*pb.GetHelloWorldResponse, error) {
	return &pb.GetHelloWorldResponse{Status: codes.OK.String(), Output: "Hello World"}, nil
}
