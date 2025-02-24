package v1

import "fmt"

// ErrValidationError is returned when invalid input is submitted for message validation.
var ErrValidationError = fmt.Errorf("failed to validate message")
