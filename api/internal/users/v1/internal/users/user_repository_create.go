package users

// create atomically checks unique constraints (email, handle, userID) and stores the entity.
func (s *userStore) create(u *User) error {
	if u == nil {
		return ErrInvalidArgument
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.users {
		if existing.EmailAddress == u.EmailAddress {
			return ErrAlreadyExists
		}
		if existing.Handle == u.Handle {
			return ErrAlreadyExists
		}
	}
	if _, ok := s.users[u.UserID]; ok {
		return ErrAlreadyExists
	}

	s.users[u.UserID] = u
	return nil
}
