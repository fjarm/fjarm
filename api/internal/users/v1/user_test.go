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
		usr      user
		contains string
		excludes string
		err      bool
	}{}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := storageUserToWireUser(&tc.usr)
			if err != nil && !tc.err {
				t.Errorf("storageUserToWireUser got an unexpected error: %v", err)
			}
			actual := redactedUserMessageString(msg)
			if !strings.Contains(actual, tc.contains) {
				t.Errorf("redactedUserMessageString got: %v, must contain: %v", actual, tc.contains)
			}
			if strings.Contains(actual, tc.excludes) {
				t.Errorf("redactedUserMessageString got: %v, must exclude: %v", actual, tc.excludes)
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
