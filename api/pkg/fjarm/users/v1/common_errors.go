package v1

import "fmt"

// ErrInputError is returned when invalid input is submitted for message validation.
var ErrInputError = fmt.Errorf("invalid input used for message validation")
