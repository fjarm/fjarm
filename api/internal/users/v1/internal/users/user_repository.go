package users

import "sync"

// userStore encapsulates thread-safe in-memory storage for user entities.
type userStore struct {
	mu    sync.RWMutex
	users map[string]*User
}

func newUserStore() *userStore {
	return &userStore{
		users: make(map[string]*User),
	}
}

func (s *userStore) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = make(map[string]*User)
}
