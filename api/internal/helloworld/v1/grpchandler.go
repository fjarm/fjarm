package v1

import (
	rpc "buf.build/gen/go/fjarm/fjarm/grpc/go/fjarm/helloworld/v1/helloworldv1grpc"
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"log/slog"
)

const handlerTag = "grpchandler"

type getHelloWorlder interface {
	getHelloWorld(ctx context.Context, input string) (string, error)
}

// GrpcHandler defines a gPRC helloworld service handler.
type GrpcHandler struct {
	rpc.UnimplementedHelloWorldServiceServer

	validator *protovalidate.Validator
	domain    getHelloWorlder
}

// NewGrpcHandler creates a gRPC handler for the helloworld service. Protovalidate is enabled by default.
func NewGrpcHandler() *GrpcHandler {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithFailFast(true),
		protovalidate.WithMessages(
			&pb.GetHelloWorldRequest{},
			&pb.GetHelloWorldResponse{},
		),
	)
	if err != nil {
		slog.Error(
			"failed to create validator",
			slog.String(logkeys.Tag, handlerTag),
			slog.String(logkeys.Err, err.Error()),
		)
		return nil
	}

	repository := newInMemoryRepository()
	domain := newInteractor(repository)

	return &GrpcHandler{
		validator: validator,
		domain:    domain,
	}
}

// GetHelloWorld implements the similarly named PRC defined in the helloworld service.
func (svc *GrpcHandler) GetHelloWorld(
	ctx context.Context,
	req *pb.GetHelloWorldRequest,
) (*pb.GetHelloWorldResponse, error) {
	slog.InfoContext(
		ctx,
		"received request",
		slog.String(logkeys.Tag, handlerTag),
		slog.String(logkeys.Rpc, "GetHelloWorld"),
		slog.String(logkeys.Request, req.String()),
	)

	err := svc.validator.Validate(req)
	if err != nil {
		slog.WarnContext(
			ctx,
			"failed to validate request",
			slog.String(logkeys.Tag, handlerTag),
			slog.String(logkeys.Rpc, "GetHelloWorld"),
			slog.Any(logkeys.Raw, req),
		)
		return nil, err
	}

	msg, err := svc.domain.getHelloWorld(ctx, req.GetInput().GetInput())
	if err != nil {
		return buildResponse(int32(code.Code_UNKNOWN), err.Error(), ""), err
	}

	return buildResponse(int32(code.Code_OK), "OK", msg), nil
}

func buildResponse(sc int32, sm string, output string) *pb.GetHelloWorldResponse {
	return &pb.GetHelloWorldResponse{
		Status: &status.Status{
			Code:    sc,
			Message: sm,
			Details: nil,
		},
		Output: &pb.HelloWorldOutput{
			Output: &output,
		},
	}
}
