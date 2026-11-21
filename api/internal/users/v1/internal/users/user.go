package users

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User defines the entity as represented in persisted storage like a SQL database.
//
// Note  that pointer fields are optional. Thus, pointers are used to distinguish between a valid value, an explicit
// null (nil pointer), or a zero value.
type User struct {
	UserID       string    `db:"user_id"`
	Handle       string    `db:"handle"`
	EmailAddress string    `db:"email_address"`
	Password     string    `db:"password"`     // Argon2ID hashed password with salt embedded.
	Avatar       *string   `db:"avatar"`       // Optional field/column. Can be nil/NULL.
	UpdatedAt    time.Time `db:"last_updated"` // Used primarily for ETag calculation.
	CreatedAt    time.Time `db:"created_at"`
}

// NewUser validates and normalizes the supplied email address, validates the hashed password, generates a random user
// handle and unique user ID, and returns a constructed User domain entity.
func NewUser(email, hashedPassword string) (*User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("%w: email address is required", ErrInvalidArgument)
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid email address", ErrInvalidArgument)
	}
	if strings.TrimSpace(hashedPassword) == "" {
		return nil, fmt.Errorf("%w: password is required", ErrInvalidArgument)
	}

	normalizedEmail := strings.ToLower(addr.Address)
	var b [16]byte
	_, err = rand.Read(b[:])
	if err != nil {
		return nil, ErrOperationFailed
	}
	username := fmt.Sprintf("user_%s", hex.EncodeToString(b[:]))

	return &User{
		UserID:       uuid.NewString(),
		Handle:       username,
		EmailAddress: normalizedEmail,
		Password:     hashedPassword,
	}, nil
}
