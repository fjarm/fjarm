package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrUnimplemented is returned when a service method is called before it's been implemented.
var ErrUnimplemented = status.Error(codes.Unimplemented, "rpc is not implemented")
