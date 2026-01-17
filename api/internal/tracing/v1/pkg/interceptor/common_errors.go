package interceptor

import (
	"fmt"

	"connectrpc.com/connect"
)

// ErrRequestIDNotFound is returned when an incoming request does not contain a `request-id` key/value pair.
var ErrRequestIDNotFound = connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to find request-id value"))
