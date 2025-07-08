package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemorySessionStore struct {
	sessions    map[string]Session
	userToToken map[string]string
	mu          sync.Mutex
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		sessions:    make(map[string]Session),
		userToToken: make(map[string]string),
	}
}

func (s *MemorySessionStore) Create(userID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldToken, exists := s.userToToken[userID]; exists {
		delete(s.sessions, oldToken)
	}

	token := uuid.New().String()
	sess := Session{
		Token:    token,
		UserID:   userID,
		LastSeen: time.Now(),
	}
	s.sessions[token] = sess
	s.userToToken[userID] = token
	return token, nil
}

func (s *MemorySessionStore) Get(token string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[token]
	if !ok {
		return nil, errors.New("not found")
	}
	return &sess, nil
}

func (s *MemorySessionStore) Delete(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[token]
	if ok {
		delete(s.sessions, token)
		delete(s.userToToken, sess.UserID)
	}
	return nil
}

func (s *MemorySessionStore) DeleteByUserID(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token, ok := s.userToToken[userID]; ok {
		delete(s.sessions, token)
		delete(s.userToToken, userID)
	}
	return nil
}

func (s *MemorySessionStore) UpdateLastSeen(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[token]
	if !ok {
		return errors.New("not found")
	}
	sess.LastSeen = time.Now()
	s.sessions[token] = sess
	return nil
}

func (s *MemorySessionStore) CleanupExpired(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, sess := range s.sessions {
		if now.Sub(sess.LastSeen) > ttl {
			delete(s.sessions, token)
			delete(s.userToToken, sess.UserID)
		}
	}
}
