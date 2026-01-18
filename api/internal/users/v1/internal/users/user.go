package users

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	consistencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/consistency/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"google.golang.org/protobuf/proto"
)

// user defines the entity as represented in persisted storage like a SQL database. This type is internal to the users
// service. External services get back a `fjarm.users.v1.User` Protobuf message instead.
//
// Note  that pointer fields are optional. Thus, pointers are used to distinguish between a valid value, an explicit
// null (nil pointer), or a zero value.
//
// TODO(2026-01-17): Remove the following fields - GivenName, FamilyName.
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

func calculateETag(usr *user) (string, error) {
	if usr == nil {
		return "", fmt.Errorf("%w: cannot calculate eTag for nil user", ErrOperationFailed)
	}

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
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func storageUserToWireUser(storageUser *user) (*userspb.User, error) {
	if storageUser == nil {
		return nil, fmt.Errorf("%w: cannot convert nil user to user message", ErrInvalidArgument)
	}
	userMsg := userspb.User{
		UserId: &userspb.UserId{UserId: &storageUser.UserID},
		FullName: &userspb.UserFullName{
			FamilyName: &storageUser.FamilyName,
			GivenName:  &storageUser.GivenName,
		},
		Handle:       &userspb.UserHandle{Handle: &storageUser.Handle},
		EmailAddress: &userspb.UserEmailAddress{EmailAddress: &storageUser.EmailAddress},
		Avatar:       &userspb.UserAvatar{Avatar: storageUser.Avatar},
	}
	etag, err := calculateETag(storageUser)
	if err != nil {
		return nil, err
	}
	userMsg.SetETag(&consistencypb.EntityTag{EntityTag: &etag})

	return &userMsg, nil
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
