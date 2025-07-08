package rbac

type Action string
type Resource string

type Permission struct {
	Action   Action
	Resource Resource
}

type Role struct {
	ID          string
	Name        string
	Permissions []Permission
}

type UserRole struct {
	UserID string
	RoleID string
}
