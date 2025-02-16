package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"time"
)

// user defines the entity as represented in persisted storage like a SQL database. This type is internal to the users
// service. External services get back a `fjarm.users.v1.User` Protobuf message instead.
type user struct {
	UserID       string    `db:"user_id"`
	GivenName    string    `db:"given_name"`
	FamilyName   string    `db:"family_name"`
	UserHandle   string    `db:"handle"`
	EmailAddress string    `db:"email_address"`
	Avatar       string    `db:"avatar"`
	Password     string    `db:"password"`
	Salt         string    `db:"salt"`
	LastUpdated  time.Time `db:"last_updated"`
	CreatedAt    time.Time `db:"created_at"`
}

func wireUserToStorageUser(msg *userspb.User) (*user, error) {
	if msg == nil {
		return nil, ErrInvalidArgument
	}
	usr := user{
		UserID:       msg.GetUserId().GetUserId(),
		GivenName:    msg.GetFullName().GetGivenName(),
		FamilyName:   msg.GetFullName().GetFamilyName(),
		UserHandle:   msg.GetHandle().GetHandle(),
		EmailAddress: msg.GetEmailAddress().GetEmailAddress(),
		Avatar:       msg.GetAvatar().GetAvatar(),
	}
	return &usr, nil
}
