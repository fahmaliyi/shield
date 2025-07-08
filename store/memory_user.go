package store

import (
	"errors"
	"sync"

	"github.com/fahmaliyi/shield/model"
)

type MemoryUserStore struct {
	usersByID    map[string]model.CoreUser
	usersByEmail map[string]model.CoreUser
	mu           sync.Mutex
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		usersByID:    make(map[string]model.CoreUser),
		usersByEmail: make(map[string]model.CoreUser),
	}
}

func (s *MemoryUserStore) Create(user model.CoreUser) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByEmail[user.GetEmail()]; exists {
		return errors.New("user already exists")
	}

	s.usersByID[user.GetID()] = user
	s.usersByEmail[user.GetEmail()] = user
	return nil
}

func (s *MemoryUserStore) FindByEmail(email string) (model.CoreUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.usersByEmail[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (s *MemoryUserStore) FindByID(userID string) (model.CoreUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.usersByID[userID]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}
