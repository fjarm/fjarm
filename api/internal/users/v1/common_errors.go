package v1

import (
	"fmt"
)

// ErrUnimplemented is returned when a service method is called before it's been implemented.
var ErrUnimplemented = fmt.Errorf("method is not implemented")

// ErrInvalidArgument is returned when a client supplies invalid arguments like a `nil` pointer to a
// `fjarm.users.v1.User` message.
var ErrInvalidArgument = fmt.Errorf("supplied argument is invalid")
