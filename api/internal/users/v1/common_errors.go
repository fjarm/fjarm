package v1

import (
	"fmt"
)

// ErrAlreadyExists is returned when a client tries to create a duplicate resource, in this case a user entity.
var ErrAlreadyExists = fmt.Errorf("resource already exists")

// ErrAuthenticationIssue is returned when a client fails to authenticate a user. This can happen at user creation when
// attempting to hash the supplied credentials. It can also fail at user update or deletion when the supplied
// credentials do not match the stored credentials.
var ErrAuthenticationIssue = fmt.Errorf("authentication issue")

// ErrInvalidArgument is returned when a client supplies invalid arguments like a `nil` pointer to a
// `fjarm.users.v1.User` message.
var ErrInvalidArgument = fmt.Errorf("supplied argument is invalid")

// ErrOperationFailed is returned when a service method fails to complete an operation. It is typically used to hide a
// more precise error message from the client like ErrAuthenticationIssue or ErrAlreadyExists.
var ErrOperationFailed = fmt.Errorf("operation failed")

// ErrUnimplemented is returned when a service method is called before it's been implemented.
var ErrUnimplemented = fmt.Errorf("method is not implemented")
