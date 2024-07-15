package helloworld

import (
	"context"
	pb "github.com/fjarm/fjarm/api/pkg/helloworld/v1"
	"google.golang.org/grpc/codes"
)

type repository struct{}

func (repo repository) getHelloWorld(_ context.Context) (*pb.GetHelloWorldResponse, error) {
	return &pb.GetHelloWorldResponse{Status: codes.Unknown.String(), Output: "Hello World"}, nil
}
