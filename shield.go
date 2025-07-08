package shield

import (
	"errors"
	"time"

	"github.com/fahmaliyi/shield/config"
	"github.com/fahmaliyi/shield/model"
	"github.com/fahmaliyi/shield/rbac"
	"github.com/fahmaliyi/shield/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrSessionExpired = errors.New("session expired")
var ErrInvalidCredentials = errors.New("invalid credentials")

type Manager struct {
	config       config.Config
	userStore    store.UserStore
	sessionStore store.SessionStore

	rbacStore rbac.RBACStore
	enforcer  *rbac.Enforcer
}

func New() *Manager {
	memRBAC := rbac.NewMemoryStore()

	return &Manager{
		config:       config.DefaultConfig(),
		userStore:    store.NewMemoryUserStore(),
		sessionStore: store.NewMemorySessionStore(),

		rbacStore: memRBAC,
		enforcer:  rbac.NewEnforcer(memRBAC),
	}
}

// Signup registers a new user
func (m *Manager) Signup(email, password string) error {
	if email == "" || password == "" {
		return errors.New("email and password required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashed),
	}

	return m.userStore.Create(user)
}

// Login verifies credentials and creates a session token
func (m *Manager) Login(email, password string) (string, error) {
	user, err := m.userStore.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.GetPasswordHash()), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	return m.sessionStore.Create(user.GetID())
}

// Logout deletes the session associated with the token
func (m *Manager) Logout(token string) error {
	return m.sessionStore.Delete(token)
}

// ValidateToken returns userID if token is valid and not expired
func (m *Manager) ValidateToken(token string) (string, error) {
	sess, err := m.sessionStore.Get(token)
	if err != nil {
		return "", err
	}

	if time.Since(sess.LastSeen) > m.config.SessionTTL {
		_ = m.sessionStore.Delete(token)
		return "", ErrSessionExpired
	}

	_ = m.sessionStore.UpdateLastSeen(token)
	return sess.UserID, nil
}

// StartCleanup starts the background session expiration cleanup goroutine
func (m *Manager) StartCleanup() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			m.sessionStore.CleanupExpired(m.config.SessionTTL)
		}
	}()
}

func (m *Manager) Can(userID string, action, resource string) bool {
	return m.enforcer.Can(userID, rbac.Action(action), rbac.Resource(resource))
}

func (m *Manager) CanOwn(userID, action, resource, ownerID string) bool {
	return m.enforcer.CanOwn(userID, rbac.Action(action), rbac.Resource(resource), ownerID)
}

func (m *Manager) AddRole(role rbac.Role) {
	m.rbacStore.AddRole(role)
}

func (m *Manager) AssignRole(userID, roleName string) {
	m.rbacStore.AssignRole(userID, roleName)
}
