package passwords

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		password string
		err      bool
	}{
		"valid_password": {
			password: "password",
			err:      false,
		},
		"valid_spaces_only_password": {
			password: "       ",
			err:      false,
		},
		"valid_non_alphabetic_password": {
			password: "$$%%$$$%%%",
			err:      false,
		},
		"valid_mixed_character_and_alphabetic_password": {
			password: "$$%%$abcd$efgh$%%%",
			err:      false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hashed, err := HashPassword(tc.password)

			if err != nil && !tc.err {
				t.Errorf("HashPassword got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Error("HashPassword expected an error but got nil")
			}

			creds, err := decodeHash(hashed)
			if err != nil {
				t.Errorf("HashPassword output cannot be decoded: %v", err)
			}

			duped, _ := base64.RawStdEncoding.DecodeString(strings.Split(hashed, delimiter)[5])
			if !bytes.Equal(duped, creds.hash) {
				t.Errorf("HashPassword output and decoded value are not equal, got: %v, want: %v", []byte(hashed), creds.hash)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	tests := map[string]struct {
		input     string
		verify    string
		match     bool
		hashErr   bool
		verifyErr bool
	}{
		"valid_matching_password": {
			input:     "password",
			verify:    "password",
			match:     true,
			hashErr:   false,
			verifyErr: false,
		},
		"valid_non_matching_password": {
			input:     "password",
			verify:    "password69",
			match:     false,
			hashErr:   false,
			verifyErr: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hashed, err := HashPassword(tc.input)
			if err != nil && !tc.hashErr {
				t.Errorf("HashPassword got an unexpected error: %v", err)
			}
			if err == nil && tc.hashErr {
				t.Error("HashPassword expected an error but got nil")
			}

			verified, err := VerifyPassword(tc.verify, hashed)
			if err != nil && !tc.verifyErr {
				t.Errorf("VerifyPassword got an unexpected error: %v", err)
			}
			if err == nil && tc.verifyErr {
				t.Error("VerifyPassword expected an error but got nil")
			}
			if verified != tc.match {
				t.Errorf("VerifyPassword got: %v, want: %v", verified, tc.match)
			}
		})
	}
}
