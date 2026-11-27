package passwords

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"testing/synctest"
)

func TestNewHasherWithMaxConcurrency(t *testing.T) {
	tests := map[string]struct {
		givenMaxConcurrency int
		wantMaxConcurrency  int
		wantErr             bool
		wantErrKind         error
	}{
		"valid_single_slot_concurrency": {
			givenMaxConcurrency: 1,
			wantMaxConcurrency:  1,
			wantErr:             false,
		},
		"valid_multi_slot_concurrency": {
			givenMaxConcurrency: 4,
			wantMaxConcurrency:  4,
			wantErr:             false,
		},
		"invalid_zero_concurrency": {
			givenMaxConcurrency: 0,
			wantErr:             true,
			wantErrKind:         ErrInvalidMaxConcurrency,
		},
		"invalid_negative_concurrency": {
			givenMaxConcurrency: -3,
			wantErr:             true,
			wantErrKind:         ErrInvalidMaxConcurrency,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			maxConcurrency := tc.givenMaxConcurrency

			// When
			hasher, err := NewHasherWithMaxConcurrency(maxConcurrency)

			// Then
			if tc.wantErr {
				if err == nil {
					t.Fatalf("NewHasherWithMaxConcurrency(%d) expected error, got nil", maxConcurrency)
				}
				if !errors.Is(err, tc.wantErrKind) {
					t.Errorf("NewHasherWithMaxConcurrency(%d) error = %v, want %v", maxConcurrency, err, tc.wantErrKind)
				}
				if hasher != nil {
					t.Errorf("NewHasherWithMaxConcurrency(%d) expected nil Hasher on error, got %v", maxConcurrency, hasher)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewHasherWithMaxConcurrency(%d) unexpected error: %v", maxConcurrency, err)
			}
			if hasher.MaxConcurrency() != tc.wantMaxConcurrency {
				t.Errorf("MaxConcurrency() = %d, want %d", hasher.MaxConcurrency(), tc.wantMaxConcurrency)
			}
		})
	}
}

func TestHasher_HashPassword(t *testing.T) {
	tests := map[string]struct {
		givenPassword       string
		givenPreCanceledCtx bool
		wantErr             bool
		wantErrKind         error
	}{
		"valid_password": {
			givenPassword: "password",
			wantErr:       false,
		},
		"valid_spaces_only_password": {
			givenPassword: "       ",
			wantErr:       false,
		},
		"valid_non_alphabetic_password": {
			givenPassword: "$$%%$$$%%%",
			wantErr:       false,
		},
		"valid_mixed_character_and_alphabetic_password": {
			givenPassword: "$$%%$abcd$efgh$%%%",
			wantErr:       false,
		},
		"canceled_context_before_hashing": {
			givenPassword:       "password",
			givenPreCanceledCtx: true,
			wantErr:             true,
			wantErrKind:         context.Canceled,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			hasher := NewHasher()
			if hasher.MaxConcurrency() < 1 {
				t.Fatalf("NewHasher().MaxConcurrency() = %d, want >= 1", hasher.MaxConcurrency())
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.givenPreCanceledCtx {
				cancel()
			}

			// When
			hashed, err := hasher.HashPassword(ctx, tc.givenPassword)

			// Then
			if tc.wantErr {
				if err == nil {
					t.Fatal("HashPassword expected an error, got nil")
				}
				if !errors.Is(err, tc.wantErrKind) {
					t.Errorf("HashPassword error = %v, want %v", err, tc.wantErrKind)
				}
				if hashed != "" {
					t.Errorf("HashPassword expected empty string on error, got %q", hashed)
				}
				return
			}

			if err != nil {
				t.Fatalf("HashPassword got an unexpected error: %v", err)
			}
			if !strings.HasPrefix(hashed, "$argon2id") {
				t.Errorf("HashPassword expected to use argon2id, got: %v", strings.Split(hashed, delimiter)[1])
			}
			if strings.Split(hashed, delimiter)[2] != "v=19" {
				t.Errorf("HashPassword expected to use argon2id version 19, got: %v", strings.Split(hashed, delimiter)[2])
			}

			creds, decodeErr := decodeHash(hashed)
			if decodeErr != nil {
				t.Fatalf("HashPassword output cannot be decoded: %v", decodeErr)
			}

			duped, b64Err := base64.RawStdEncoding.DecodeString(strings.Split(hashed, delimiter)[5])
			if b64Err != nil {
				t.Fatalf("failed to decode base64 hash segment: %v", b64Err)
			}
			if !bytes.Equal(duped, creds.hash) {
				t.Errorf("HashPassword output and decoded value are not equal, got: %v, want: %v", duped, creds.hash)
			}
		})
	}
}

func TestHasher_VerifyPassword(t *testing.T) {
	tests := map[string]struct {
		givenInput           string
		givenVerify          string
		givenEncodedOverride string
		givenPreCanceledCtx  bool
		wantMatch            bool
		wantVerifyErr        bool
		wantErrKind          error
	}{
		"valid_matching_password": {
			givenInput:    "password",
			givenVerify:   "password",
			wantMatch:     true,
			wantVerifyErr: false,
		},
		"valid_non_matching_password": {
			givenInput:    "password",
			givenVerify:   "password69",
			wantMatch:     false,
			wantVerifyErr: false,
		},
		"invalid_encoded_hash_format": {
			givenVerify:          "password",
			givenEncodedOverride: "invalid-format",
			wantMatch:            false,
			wantVerifyErr:        true,
			wantErrKind:          ErrInvalidHashFormat,
		},
		"canceled_context_before_verifying": {
			givenInput:          "password",
			givenVerify:         "password",
			givenPreCanceledCtx: true,
			wantMatch:           false,
			wantVerifyErr:       true,
			wantErrKind:         context.Canceled,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			hasher := NewHasher()

			encoded := tc.givenEncodedOverride
			if encoded == "" {
				var hashErr error
				encoded, hashErr = hasher.HashPassword(context.Background(), tc.givenInput)
				if hashErr != nil {
					t.Fatalf("setup HashPassword got an unexpected error: %v", hashErr)
				}
			}

			verifyCtx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.givenPreCanceledCtx {
				cancel()
			}

			// When
			verified, err := hasher.VerifyPassword(verifyCtx, tc.givenVerify, encoded)

			// Then
			if tc.wantVerifyErr {
				if err == nil {
					t.Fatal("VerifyPassword expected an error, got nil")
				}
				if !errors.Is(err, tc.wantErrKind) {
					t.Errorf("VerifyPassword error = %v, want %v", err, tc.wantErrKind)
				}
			} else if err != nil {
				t.Errorf("VerifyPassword got an unexpected error: %v", err)
			}

			if verified != tc.wantMatch {
				t.Errorf("VerifyPassword got: %v, want: %v", verified, tc.wantMatch)
			}
		})
	}
}

func TestHasher_SemaphoreConcurrencyControl(t *testing.T) {
	fastParams := &hashParams{
		memory:      64,
		iterations:  1,
		parallelism: 1,
		saltLength:  16,
		keyLength:   32,
	}

	setupHasher, setupErr := NewHasherWithMaxConcurrency(1)
	if setupErr != nil {
		t.Fatalf("setup NewHasherWithMaxConcurrency got unexpected error: %v", setupErr)
	}
	setupHasher.params = fastParams

	encodedCredential, encodeErr := setupHasher.HashPassword(context.Background(), "password")
	if encodeErr != nil {
		t.Fatalf("setup HashPassword got unexpected error: %v", encodeErr)
	}

	tests := map[string]struct {
		givenMaxConcurrency int
		givenPreFilledSlots int
		givenOperation      string
		givenCancelQueued   bool
		wantBlockedAtStart  bool
		wantErr             bool
		wantErrKind         error
	}{
		"hash_blocks_when_all_slots_occupied_then_succeeds_after_release": {
			givenMaxConcurrency: 1,
			givenPreFilledSlots: 1,
			givenOperation:      "hash",
			givenCancelQueued:   false,
			wantBlockedAtStart:  true,
			wantErr:             false,
		},
		"verify_blocks_when_all_slots_occupied_then_succeeds_after_release": {
			givenMaxConcurrency: 1,
			givenPreFilledSlots: 1,
			givenOperation:      "verify",
			givenCancelQueued:   false,
			wantBlockedAtStart:  true,
			wantErr:             false,
		},
		"hash_proceeds_immediately_when_spare_slot_is_available": {
			givenMaxConcurrency: 2,
			givenPreFilledSlots: 1,
			givenOperation:      "hash",
			givenCancelQueued:   false,
			wantBlockedAtStart:  false,
			wantErr:             false,
		},
		"hash_aborts_on_context_cancel_when_all_slots_occupied": {
			givenMaxConcurrency: 1,
			givenPreFilledSlots: 1,
			givenOperation:      "hash",
			givenCancelQueued:   true,
			wantBlockedAtStart:  true,
			wantErr:             true,
			wantErrKind:         context.Canceled,
		},
		"verify_aborts_on_context_cancel_when_all_slots_occupied": {
			givenMaxConcurrency: 1,
			givenPreFilledSlots: 1,
			givenOperation:      "verify",
			givenCancelQueued:   true,
			wantBlockedAtStart:  true,
			wantErr:             true,
			wantErrKind:         context.Canceled,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				// Given
				hasher, err := NewHasherWithMaxConcurrency(tc.givenMaxConcurrency)
				if err != nil {
					t.Fatalf("NewHasherWithMaxConcurrency got unexpected error: %v", err)
				}
				hasher.params = fastParams

				for range tc.givenPreFilledSlots {
					hasher.semaphore <- struct{}{}
				}

				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()

				// When
				errCh := make(chan error, 1)
				go func() {
					if tc.givenOperation == "verify" {
						_, opErr := hasher.VerifyPassword(ctx, "password", encodedCredential)
						errCh <- opErr
						return
					}
					_, opErr := hasher.HashPassword(ctx, "password")
					errCh <- opErr
				}()

				synctest.Wait()

				var gotErr error
				remainingPreFilled := tc.givenPreFilledSlots

				if tc.wantBlockedAtStart {
					select {
					case prematureErr := <-errCh:
						t.Fatalf("operation completed while semaphore was full (err: %v)", prematureErr)
					default:
					}

					if tc.givenCancelQueued {
						cancel()
						synctest.Wait()
						gotErr = <-errCh
					} else {
						<-hasher.semaphore
						remainingPreFilled--
						synctest.Wait()
						gotErr = <-errCh
					}
				} else {
					select {
					case gotErr = <-errCh:
					default:
						t.Fatal("operation blocked despite spare semaphore slot")
					}
				}

				for range remainingPreFilled {
					<-hasher.semaphore
				}

				// Then
				if tc.wantErr {
					if gotErr == nil {
						t.Fatal("expected error, got nil")
					}
					if !errors.Is(gotErr, tc.wantErrKind) {
						t.Errorf("got error %v, want %v", gotErr, tc.wantErrKind)
					}
				} else if gotErr != nil {
					t.Errorf("got unexpected error: %v", gotErr)
				}

				if len(hasher.semaphore) != 0 {
					t.Errorf("semaphore len = %d after cleanup, want 0", len(hasher.semaphore))
				}
			})
		})
	}
}
