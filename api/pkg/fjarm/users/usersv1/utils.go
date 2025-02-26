package usersv1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"strings"
)

func ValidateUserEmailAddress(_ context.Context, email *userspb.UserEmailAddress) error {
	if email == nil {
		return ErrValidationError
	}
	if !email.HasEmailAddress() {
		return ErrValidationError
	}
	return protovalidate.Validate(email)
}

func ValidateUserFullName(_ context.Context, name *userspb.UserFullName) error {
	if name == nil {
		return ErrValidationError
	}
	if !name.HasGivenName() || !name.HasFamilyName() {
		return ErrValidationError
	}
	return protovalidate.Validate(name)
}

func ValidateUserHandle(_ context.Context, handle *userspb.UserHandle) error {
	if handle == nil {
		return ErrValidationError
	}
	if !handle.HasHandle() {
		return ErrValidationError
	}
	if strings.Contains(handle.GetHandle(), " ") {
		return ErrValidationError
	}
	if handle.GetHandle() == "" {
		return ErrValidationError
	}
	return protovalidate.Validate(handle)

}

func ValidateUserID(_ context.Context, id *userspb.UserId) error {
	if id == nil {
		return ErrValidationError
	}
	if !id.HasUserId() {
		return ErrValidationError
	}
	return protovalidate.Validate(id)
}

func ValidateUserPassword(_ context.Context, pwd *userspb.UserPassword) error {
	if pwd == nil {
		return ErrValidationError
	}
	if !pwd.HasPassword() {
		return ErrValidationError
	}
	return protovalidate.Validate(pwd)
}
