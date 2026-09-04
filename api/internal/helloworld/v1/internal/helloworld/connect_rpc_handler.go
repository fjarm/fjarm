package helloworld

import (
	"context"
	"log/slog"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/helloworld/v1/helloworldv1connect"
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"buf.build/go/protovalidate"
	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const connectRPCHandlerTag = "connect_rpc_handler"

type getHelloWorlder interface {
	getHelloWorld(ctx context.Context, input string) (string, error)
}

// ConnectRPCHandler defines a ConnectRPC handler for the `fjarm.helloworld.v1.HelloWorldService` service.
type ConnectRPCHandler struct {
	helloworldv1connect.UnimplementedHelloWorldServiceHandler
	domain    getHelloWorlder
	logger    *slog.Logger
	validator protovalidate.Validator
}

// GetHelloWorld implements the similarly named RPC defined in the `fjarm.helloworld.v1.HelloWorldService` service.
func (h *ConnectRPCHandler) GetHelloWorld(
	ctx context.Context,
	req *pb.GetHelloWorldRequest,
) (*pb.GetHelloWorldResponse, error) {
	requestID := tracing.RequestIDFromContext(ctx)

	logger := h.logger.With(
		slog.String(logkeys.Rpc, helloworldv1connect.HelloWorldServiceGetHelloWorldProcedure),
		slog.String(tracing.RequestIDKey, requestID),
	)

	logger.InfoContext(
		ctx,
		"received request",
		slog.String(logkeys.Request, req.String()),
	)

	err := h.validator.Validate(req)
	if err != nil {
		logger.WarnContext(ctx, "failed to validate request", slog.Any(logkeys.Raw, req))
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	msg, err := h.domain.getHelloWorld(ctx, req.GetInput().GetInput())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := &pb.GetHelloWorldResponse{
		Output: &pb.HelloWorldOutput{
			Output: &msg,
		},
	}
	return res, nil
}

// NewConnectRPCHandler creates a concrete helloworld
func NewConnectRPCHandler(l *slog.Logger) *ConnectRPCHandler {
	logger := l.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
	)

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&pb.GetHelloWorldRequest{},
			&pb.GetHelloWorldResponse{},
		),
	)
	if err != nil {
		logger.Error("failed to create message validator", slog.Any(logkeys.Err, err))
		return nil
	}

	repo := newInMemoryRepository()
	dom := newInteractor(logger, repo)
	han := ConnectRPCHandler{
		domain:    dom,
		logger:    logger,
		validator: validator,
	}
	return &han
}
