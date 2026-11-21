package users

import (
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
			wantPrefix:     "foo",
			wantErr:        false,
		},
		"valid_email_with_display_name_is_normalized": {
			email:          "Foo Bar <foo@bar.com>",
			hashedPassword: "hashed_password_123",
			wantEmail:      "foo@bar.com",
			wantPrefix:     "foo",
			wantErr:        false,
		},
		"valid_email_with_mixed_case_is_normalized": {
			email:          "Foo@Bar.com",
			hashedPassword: "hashed_password_123",
			wantEmail:      "foo@bar.com",
			wantPrefix:     "foo",
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
			if !strings.HasPrefix(usr.Handle, tc.wantPrefix) || len(usr.Handle) <= len(tc.wantPrefix) {
				t.Errorf("NewUser Handle = %q, expected prefix %q followed by FNV-64a hash", usr.Handle, tc.wantPrefix)
			}
		})
	}
}

func TestNewUser_DeterministicHandleAndUniqueID(t *testing.T) {
	u1, err := NewUser("foo@bar.com", "hashed_pwd_1")
	if err != nil {
		t.Fatalf("NewUser unexpected error: %v", err)
	}
	u2, err := NewUser("Foo <Foo@Bar.com>", "hashed_pwd_2")
	if err != nil {
		t.Fatalf("NewUser unexpected error: %v", err)
	}

	if u1.Handle != u2.Handle {
		t.Errorf("expected identical handles for same normalized email, got %q and %q", u1.Handle, u2.Handle)
	}
	if u1.UserID == u2.UserID {
		t.Errorf("expected distinct UserIDs across NewUser calls, got %q", u1.UserID)
	}
}
