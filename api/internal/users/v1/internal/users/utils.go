package users

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"github.com/fjarm/fjarm/api/internal/logvals"
	"github.com/fjarm/fjarm/api/pkg/fjarm/users/usersv1"
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

func validateUserMessageForCreate(ctx context.Context, msg *userspb.User) error {
	if msg == nil {
		return ErrInvalidArgument
	}

	err := usersv1.ValidateUserID(ctx, msg.GetUserId())
	if err != nil {
		return err
	}

	err = usersv1.ValidateUserHandle(ctx, msg.GetHandle())
	if err != nil {
		return err
	}

	err = usersv1.ValidateUserEmailAddress(ctx, msg.GetEmailAddress())
	if err != nil {
		return err
	}

	err = usersv1.ValidateUserFullName(ctx, msg.GetFullName())
	if err != nil {
		return err
	}

	err = usersv1.ValidateUserPassword(ctx, msg.GetPassword())
	if err != nil {
		return err
	}

	return nil
}
