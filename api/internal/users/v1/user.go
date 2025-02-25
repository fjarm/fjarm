package v1

import (
	consistencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/consistency/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"time"
)

// user defines the entity as represented in persisted storage like a SQL database. This type is internal to the users
// service. External services get back a `fjarm.users.v1.User` Protobuf message instead.
//
// Note  that pointer fields are optional. Thus, pointers are used to distinguish between a valid value, an explicit
// null (nil pointer), or a zero value.
type user struct {
	UserID       string    `db:"user_id"`
	GivenName    string    `db:"given_name"`
	FamilyName   string    `db:"family_name"`
	Handle       string    `db:"handle"`
	EmailAddress string    `db:"email_address"`
	Password     string    `db:"password"`     // Argon2ID hashed password with salt embedded.
	Avatar       *string   `db:"avatar"`       // Optional field/column. Can be nil/NULL.
	LastUpdated  time.Time `db:"last_updated"` // Used primarily for ETag calculation.
	CreatedAt    time.Time `db:"created_at"`
}

func (usr *user) calculateETag() string {
	hasher := sha256.New()

	// Note: We're explicitly ordering fields to ensure consistency
	hasher.Write([]byte(usr.UserID))
	hasher.Write([]byte(usr.GivenName))
	hasher.Write([]byte(usr.FamilyName))
	hasher.Write([]byte(usr.Handle))
	hasher.Write([]byte(usr.EmailAddress))
	if usr.Avatar != nil {
		hasher.Write([]byte(*usr.Avatar))
	}
	hasher.Write([]byte(usr.LastUpdated.UTC().Format(time.RFC3339Nano)))
	hasher.Write([]byte(usr.CreatedAt.UTC().Format(time.RFC3339Nano)))

	// Convert hash to base64 for a shorter string representation
	// Use RawURLEncoding to ensure URL-safe characters
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
}

func storageUserToWireUser(su *user) (*userspb.User, error) {
	if su == nil {
		return nil, fmt.Errorf("%w: cannot convert nil user to user message", ErrInvalidArgument)
	}
	usr := userspb.User{
		UserId: &userspb.UserId{UserId: &su.UserID},
		FullName: &userspb.UserFullName{
			FamilyName: &su.FamilyName,
			GivenName:  &su.GivenName,
		},
		Handle:       &userspb.UserHandle{Handle: &su.Handle},
		EmailAddress: &userspb.UserEmailAddress{EmailAddress: &su.EmailAddress},
		Avatar:       &userspb.UserAvatar{Avatar: su.Avatar},
	}
	etag := su.calculateETag()
	usr.SetETag(&consistencypb.EntityTag{EntityTag: &etag})

	return &usr, nil
}

func wireUserToStorageUser(msg *userspb.User) (*user, error) {
	if msg == nil {
		return nil, fmt.Errorf("%w: cannot convert nil user message to user", ErrInvalidArgument)
	}
	usr := user{
		UserID:       msg.GetUserId().GetUserId(),
		GivenName:    msg.GetFullName().GetGivenName(),
		FamilyName:   msg.GetFullName().GetFamilyName(),
		Handle:       msg.GetHandle().GetHandle(),
		EmailAddress: msg.GetEmailAddress().GetEmailAddress(),
		Avatar:       proto.String(msg.GetAvatar().GetAvatar()),
	}
	return &usr, nil
}
