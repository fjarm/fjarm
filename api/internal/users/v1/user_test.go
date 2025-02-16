package v1

import (
	"errors"
	"testing"
)

func TestWireUserToStorageUser_ErrorWrapping(t *testing.T) {
	_, err := wireUserToStorageUser(nil)
	if !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("error is not wrapped correctly")
	}
}
