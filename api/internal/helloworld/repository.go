package helloworld

import (
	"context"
	pb "github.com/fjarm/fjarm/api/pkg/helloworld/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Repository struct {
}

func (repo *Repository) GetHelloWorld(
	ctx context.Context,
	req *emptypb.Empty,
) (*pb.GetHelloWorldResponse, error) {
	return &pb.GetHelloWorldResponse{
		Status: codes.OK.String(),
		Output: "Hello world!",
	}, nil
}

func NewRepository() *Repository {
	return &Repository{}
}
