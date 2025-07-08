package store

import (
	"time"

	"github.com/fahmaliyi/shield/model"
)

type Session struct {
	Token    string
	UserID   string
	LastSeen time.Time
}

type SessionStore interface {
	Create(userID string) (string, error)
	Get(token string) (*Session, error)
	Delete(token string) error
	DeleteByUserID(userID string) error
	UpdateLastSeen(token string) error
	CleanupExpired(ttl time.Duration)
}

type UserStore interface {
	Create(user model.CoreUser) error
	FindByEmail(email string) (model.CoreUser, error)
	FindByID(userID string) (model.CoreUser, error)
}
