package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
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
