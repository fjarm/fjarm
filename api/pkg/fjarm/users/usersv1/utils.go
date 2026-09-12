package usersv1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
)

func ValidateUserEmailAddress(email *userspb.UserEmailAddress) error {
	if email == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(email)
}

func ValidateUserHandle(handle *userspb.UserHandle) error {
	if handle == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(handle)
}

func ValidateUserID(id *userspb.UserId) error {
	if id == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(id)
}

func ValidateUserPassword(pwd *userspb.UserPassword) error {
	if pwd == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(pwd)
}
