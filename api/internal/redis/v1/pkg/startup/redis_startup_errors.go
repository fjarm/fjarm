package startup

import "fmt"

// ErrUnimplemented is returned from methods/functions that haven't been implemented yet.
var ErrUnimplemented = fmt.Errorf("not implemented")
