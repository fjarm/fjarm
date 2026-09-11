package users

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"

	"github.com/fjarm/fjarm/api/internal/logvals"
)

func redactedUserMessageString(msg *userspb.User) string {
	if msg == nil {
		return logvals.Nil
	}
	rm := &userspb.User{
		UserId: msg.UserId,
	}
	return rm.String()
}

func validateUserMessageForCreate(msg *userspb.User) error {
	if msg == nil {
		return ErrInvalidArgument
	}
	// Context-specific: Create requires these sub-messages to be present
	if !msg.HasHandle() || !msg.HasEmailAddress() || !msg.HasPassword() {
		return ErrInvalidArgument
	}

	return protovalidate.Validate(msg)
}
