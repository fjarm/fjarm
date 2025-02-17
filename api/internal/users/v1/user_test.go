package v1

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWireUserToStorageUser_ErrorWrapping(t *testing.T) {
	_, err := wireUserToStorageUser(nil)
	if !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("error is not wrapped correctly")
	}
}

func TestStorageUserToWireUser_ErrorWrapping(t *testing.T) {
	_, err := storageUserToWireUser(nil)
	if !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("error is not wrapped correctly")
	}
}

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
