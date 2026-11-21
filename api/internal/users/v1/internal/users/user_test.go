package users

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tests := map[string]struct {
		email          string
		hashedPassword string
		wantEmail      string
		wantPrefix     string
		wantErr        bool
		wantErrKind    error
	}{
		"valid_email_and_password": {
			email:          "foo@bar.com",
			hashedPassword: "hashed_password_123",
			wantEmail:      "foo@bar.com",
			wantPrefix:     "user_",
			wantErr:        false,
		},
		"valid_email_with_display_name_is_normalized": {
			email:          "Foo Bar <foo@bar.com>",
			hashedPassword: "hashed_password_123",
			wantEmail:      "foo@bar.com",
			wantPrefix:     "user_",
			wantErr:        false,
		},
		"valid_email_with_mixed_case_is_normalized": {
			email:          "Foo@Bar.com",
			hashedPassword: "hashed_password_123",
			wantEmail:      "foo@bar.com",
			wantPrefix:     "user_",
			wantErr:        false,
		},
		"invalid_empty_email": {
			email:          "",
			hashedPassword: "hashed_password_123",
			wantErr:        true,
			wantErrKind:    ErrInvalidArgument,
		},
		"invalid_whitespace_email": {
			email:          "   ",
			hashedPassword: "hashed_password_123",
			wantErr:        true,
			wantErrKind:    ErrInvalidArgument,
		},
		"invalid_malformed_email": {
			email:          "not-an-email",
			hashedPassword: "hashed_password_123",
			wantErr:        true,
			wantErrKind:    ErrInvalidArgument,
		},
		"invalid_empty_password": {
			email:          "foo@bar.com",
			hashedPassword: "",
			wantErr:        true,
			wantErrKind:    ErrInvalidArgument,
		},
		"invalid_whitespace_password": {
			email:          "foo@bar.com",
			hashedPassword: "   ",
			wantErr:        true,
			wantErrKind:    ErrInvalidArgument,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			usr, err := NewUser(tc.email, tc.hashedPassword)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("NewUser(%q, %q) expected error, got nil", tc.email, tc.hashedPassword)
				}
				if !errors.Is(err, tc.wantErrKind) {
					t.Errorf("NewUser(%q, %q) error = %v, want error wrapping %v", tc.email, tc.hashedPassword, err, tc.wantErrKind)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewUser(%q, %q) unexpected error: %v", tc.email, tc.hashedPassword, err)
			}
			if _, parseErr := uuid.Parse(usr.UserID); parseErr != nil {
				t.Errorf("NewUser UserID %q is not a valid UUID: %v", usr.UserID, parseErr)
			}
			if usr.EmailAddress != tc.wantEmail {
				t.Errorf("NewUser EmailAddress = %q, want %q", usr.EmailAddress, tc.wantEmail)
			}
			if usr.Password != tc.hashedPassword {
				t.Errorf("NewUser Password = %q, want %q", usr.Password, tc.hashedPassword)
			}
			if !strings.HasPrefix(usr.Handle, tc.wantPrefix) {
				t.Errorf("NewUser Handle = %q, expected prefix %q", usr.Handle, tc.wantPrefix)
			}
			suffix := strings.TrimPrefix(usr.Handle, tc.wantPrefix)
			if decoded, hexErr := hex.DecodeString(suffix); hexErr != nil || len(decoded) != 16 {
				t.Errorf("NewUser Handle = %q, expected 16-byte hex suffix, got %q (err: %v)", usr.Handle, suffix, hexErr)
			}
		})
	}
}

func TestNewUser_UniqueHandleAndUniqueID(t *testing.T) {
	tests := map[string]struct {
		emails          []string
		hashedPasswords []string
		wantUniqueCount int
	}{
		"valid_identical_normalized_emails_produce_unique_handles_and_ids": {
			emails:          []string{"foo@bar.com", "Foo <Foo@Bar.com>"},
			hashedPasswords: []string{"hashed_pwd_1", "hashed_pwd_2"},
			wantUniqueCount: 2,
		},
		"valid_distinct_emails_produce_unique_handles_and_ids": {
			emails:          []string{"foo1@bar.com", "foo2@bar.com"},
			hashedPasswords: []string{"hashed_pwd_1", "hashed_pwd_2"},
			wantUniqueCount: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			seenHandles := make(map[string]struct{})
			seenUserIDs := make(map[string]struct{})

			for i := range tc.emails {
				usr, err := NewUser(tc.emails[i], tc.hashedPasswords[i])
				if err != nil {
					t.Fatalf("NewUser(%q, %q) unexpected error: %v", tc.emails[i], tc.hashedPasswords[i], err)
				}
				seenHandles[usr.Handle] = struct{}{}
				seenUserIDs[usr.UserID] = struct{}{}
			}

			if len(seenHandles) != tc.wantUniqueCount {
				t.Errorf("NewUser unique handles got: %d, want: %d", len(seenHandles), tc.wantUniqueCount)
			}
			if len(seenUserIDs) != tc.wantUniqueCount {
				t.Errorf("NewUser unique UserIDs got: %d, want: %d", len(seenUserIDs), tc.wantUniqueCount)
			}
		})
	}
}
