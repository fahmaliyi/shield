package rbac

import "sync"

type MemoryStore struct {
	roles      map[string]Role
	userRoles  map[string]string
	roleByName map[string]string
	mu         sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		roles:      make(map[string]Role),
		userRoles:  make(map[string]string),
		roleByName: make(map[string]string),
	}
}

func (s *MemoryStore) AddRole(role Role) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.roles[role.ID] = role
	s.roleByName[role.Name] = role.ID
}

func (s *MemoryStore) AssignRole(userID, roleName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if roleID, ok := s.roleByName[roleName]; ok {
		s.userRoles[userID] = roleID
	}
}

func (s *MemoryStore) GetPermissions(userID string) []Permission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	roleID, ok := s.userRoles[userID]
	if !ok {
		return nil
	}
	role, ok := s.roles[roleID]
	if !ok {
		return nil
	}
	return role.Permissions
}
