package users

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"

	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"google.golang.org/protobuf/proto"
)

func TestInMemoryRepository_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newInMemoryRepository(logger)

	tests := map[string]struct {
		users []*userspb.User
		err   []bool
		kind  []error
	}{
		"validation_one_valid_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false},
			kind: []error{nil},
		},
		"validation_one_nil_user": {
			users: []*userspb.User{
				nil,
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_empty_user": {
			users: []*userspb.User{
				{},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_id_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_id_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("user_id")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_empty_string_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_contains_spaces_handle_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String(" ")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_email_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_email_user": {
			users: []*userspb.User{
				{
					UserId:   &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					Handle:   &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password: &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_invalid_email_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("gleeper")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_no_password_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"validation_one_unset_password_user": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
				},
			},
			err:  []bool{true},
			kind: []error{ErrInvalidArgument},
		},
		"idempotency_two_distinct_valid_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4568-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false, false},
			kind: []error{nil, nil},
		},
		"idempotency_two_identical_id_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
		"idempotency_two_identical_email_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper1")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
		"idempotency_two_identical_handle_users": {
			users: []*userspb.User{
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo1@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
				{
					UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174999")},
					EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("foo@bar.com")},
					Handle:       &userspb.UserHandle{Handle: proto.String("gleeper")},
					Password:     &userspb.UserPassword{Password: proto.String("password")},
				},
			},
			err:  []bool{false, true},
			kind: []error{nil, ErrAlreadyExists},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for index, create := range tc.users {
				_, err := repo.createUser(context.Background(), create)
				if err != nil && !tc.err[index] {
					t.Errorf("createUser got an unexpected error: %v", err)
				}
				if err == nil && tc.err[index] {
					t.Errorf("createUser expected an error but got nil")
				}
				if !errors.Is(err, tc.kind[index]) {
					t.Errorf("createUser got an unexpected error type: %v", err)
				}
			}
		})
		// Reset the database for each test run.
		repo.reset()
	}
}

func TestInMemoryRepository_createUser_concurrency(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newInMemoryRepository(logger)

	const numGoroutines = 50
	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines)

	for i := range numGoroutines {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			user := &userspb.User{
				UserId:       &userspb.UserId{UserId: proto.String(fmt.Sprintf("123e4567-e89b-12d3-a456-426614174%03d", idx))},
				EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String(fmt.Sprintf("user%d@example.com", idx))},
				Handle:       &userspb.UserHandle{Handle: proto.String(fmt.Sprintf("user%d", idx))},
				Password:     &userspb.UserPassword{Password: proto.String("password123")},
			}
			created, err := repo.createUser(context.Background(), user)
			if err != nil {
				errChan <- err
				return
			}
			if created == nil || created.UserID == "" {
				errChan <- errors.New("created user or UserID is empty")
				return
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrent createUser failed: %v", err)
	}
}

func TestInMemoryRepository_createUser_concurrentDuplicates(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newInMemoryRepository(logger)

	const numGoroutines = 20
	var wg sync.WaitGroup

	successCount := 0
	var mu sync.Mutex

	for range numGoroutines {
		wg.Go(func() {
			user := &userspb.User{
				UserId:       &userspb.UserId{UserId: proto.String("123e4567-e89b-12d3-a456-426614174000")},
				EmailAddress: &userspb.UserEmailAddress{EmailAddress: proto.String("duplicate@example.com")},
				Handle:       &userspb.UserHandle{Handle: proto.String("duplicate")},
				Password:     &userspb.UserPassword{Password: proto.String("password123")},
			}
			_, err := repo.createUser(context.Background(), user)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			} else if !errors.Is(err, ErrAlreadyExists) {
				t.Errorf("unexpected error for duplicate creation: %v", err)
			}
		})
	}

	wg.Wait()

	if successCount != 1 {
		t.Errorf("expected exactly 1 successful creation, got %d", successCount)
	}
}
