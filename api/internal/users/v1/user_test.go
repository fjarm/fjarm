package v1

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestCalculateETag(t *testing.T) {
	tests := map[string]struct {
		u1    user
		u2    user
		equal bool
	}{
		"valid_two_empty_users": {
			u1:    user{},
			u2:    user{},
			equal: true,
		},
		"valid_two_non_empty_users": {
			u1: user{
				UserID:      "abc123",
				LastUpdated: time.Date(2024, time.December, 12, 12, 12, 12, 12, time.Local),
				CreatedAt:   time.Date(2024, time.December, 12, 12, 12, 7, 12, time.Local),
			},
			u2: user{
				UserID:      "abc123",
				LastUpdated: time.Date(2024, time.December, 12, 12, 12, 12, 12, time.Local),
				CreatedAt:   time.Date(2024, time.December, 12, 12, 12, 7, 12, time.Local),
			},
			equal: true,
		},
		"valid_different_passwords_similar_users": {
			u1: user{
				UserID:   "abc123",
				Password: "password",
			},
			u2: user{
				UserID:   "abc123",
				Password: "password1",
			},
			equal: true,
		},
		"valid_different_salts_similar_users": {
			u1: user{
				UserID: "abc123",
				Salt:   "salt",
			},
			u2: user{
				UserID: "abc123",
				Salt:   "salt2",
			},
			equal: true,
		},
		"valid_nil_and_valid_passwords_similar_users": {
			u1: user{
				UserID:   "abc123",
				Password: "password",
			},
			u2: user{
				UserID: "abc123",
			},
			equal: true,
		},
		"valid_nil_and_valid_salts_similar_users": {
			u1: user{
				UserID: "abc123",
				Salt:   "password",
			},
			u2: user{
				UserID: "abc123",
			},
			equal: true,
		},
		"invalid_two_non_empty_users": {
			u1: user{
				UserID:      "abc123",
				LastUpdated: time.Now(),
				CreatedAt:   time.Now(),
			},
			u2: user{
				UserID:      "abc123",
				LastUpdated: time.Now(),
				CreatedAt:   time.Now(),
			},
			equal: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			e1 := tc.u1.calculateETag()
			e2 := tc.u2.calculateETag()
			if e1 != e2 && tc.equal {
				t.Errorf("unexpected two unequal users")
			}
		})
	}
}

func TestRedactedUserMessageString(t *testing.T) {
	tests := map[string]struct {
		usr      *user
		contains []string
		excludes []string
		err      bool
	}{
		"valid_empty_user": {
			usr:      &user{},
			contains: []string{"{", "}", "user_id"},
			excludes: []string{"password", "nil", "salt", "full_name", "handle", "email_address", "avatar", "last_updated", "created_at"},
			err:      false,
		},
		"valid_nil_user": {
			contains: []string{"nil"},
			excludes: []string{"password", "salt", "{", "}", "user_id", "full_name", "handle", "email_address", "avatar"},
			err:      true,
		},
		"valid_valid_password_and_salt_user": {
			usr: &user{
				Password: "abc123",
				Salt:     "salt",
			},
			contains: []string{"{", "}", "user_id"},
			excludes: []string{"password", "salt", "full_name", "handle", "email_address", "avatar", "abc123", "salt"},
			err:      false,
		},
		"valid_valid_user_id_user": {
			usr: &user{
				UserID: "abc123",
				Salt:   "salt",
			},
			contains: []string{"{", "}", "user_id", "abc123"},
			excludes: []string{"password", "salt", "full_name", "handle", "email_address", "avatar"},
			err:      false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := storageUserToWireUser(tc.usr)
			if !tc.err && err != nil {
				t.Errorf("storageUserToWireUser got an unexpected error: %v", err)
			}
			actual := redactedUserMessageString(msg)
			for _, s := range tc.contains {
				if !strings.Contains(actual, s) {
					t.Errorf("redactedUserMessageString got: %v, must contain: %v", actual, tc.contains)
				}
			}
			for _, s := range tc.excludes {
				if strings.Contains(actual, s) {
					t.Errorf("redactedUserMessageString got: %v, must exclude: %v", actual, tc.excludes)
				}
			}
		})
	}
}

func TestStorageUserToWireUser_ErrorWrapping(t *testing.T) {
	_, err := storageUserToWireUser(nil)
	if !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("error is not wrapped correctly")
	}
}

func TestStorageUserToWireUser_EtagCalculation(t *testing.T) {
	tests := map[string]struct {
		usr user
		err bool
	}{
		"valid_empty_user": {
			usr: user{},
			err: false,
		},
		"valid_non_empty_user": {
			usr: user{
				UserID: "abc123",
			},
			err: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := storageUserToWireUser(&tc.usr)
			if err != nil && !tc.err {
				t.Errorf("storageUserToWireUser got an unexpected error: %v", err)
			}
			if msg.GetETag().GetEntityTag() != tc.usr.calculateETag() {
				t.Errorf(
					"unexpected etag calculation got %v, want %v",
					msg.GetETag().GetEntityTag(),
					tc.usr.calculateETag(),
				)
			}
		})
	}
}

func TestWireUserToStorageUser_ErrorWrapping(t *testing.T) {
	_, err := wireUserToStorageUser(nil)
	if !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("error is not wrapped correctly")
	}
}
