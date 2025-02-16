package v1

import (
	"connectrpc.com/connect"
	"fmt"
)

// ErrUnimplemented is returned when a service method is called before it's been implemented.
var ErrUnimplemented = connect.NewError(connect.CodeUnimplemented, fmt.Errorf("rpc is not implemented"))

// ErrInvalidArgument is returned when a client supplies invalid arguments like a `nil` pointer to a
// `fjarm.users.v1.User` message.
var ErrInvalidArgument = connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("supplied argument is invalid"))
