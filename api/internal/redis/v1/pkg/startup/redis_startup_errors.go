package startup

import "fmt"

// ErrInvalidCommand is returned when trying to construct a command permission with invalid input.
var ErrInvalidCommand = fmt.Errorf("invalid command")

// ErrUnimplemented is returned from methods/functions that haven't been implemented yet.
var ErrUnimplemented = fmt.Errorf("not implemented")
