package rbac

type RBACStore interface {
	AddRole(role Role)
	AssignRole(userID, roleName string)
	GetPermissions(userID string) []Permission
}
