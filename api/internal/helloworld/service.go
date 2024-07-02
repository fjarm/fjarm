package helloworld

import (
	"context"
	pb "github.com/fjarm/fjarm/api/pkg/helloworld/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pb.UnimplementedHelloWorldServiceServer
	repo repository
}

func (svc *Service) GetHelloWorld(
	ctx context.Context,
	req *emptypb.Empty,
) (*pb.GetHelloWorldResponse, error) {
	return svc.repo.getHelloWorld(ctx)
}

func NewService() *Service {
	return &Service{
		repo: repository{},
	}
}
