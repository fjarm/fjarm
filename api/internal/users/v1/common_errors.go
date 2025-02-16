package v1

import (
	"connectrpc.com/connect"
	"fmt"
)

// ErrUnimplemented is returned when a service method is called before it's been implemented.
var ErrUnimplemented = connect.NewError(connect.CodeUnimplemented, fmt.Errorf("rpc is not implemented"))
