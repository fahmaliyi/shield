package rbac

type Enforcer struct {
	store RBACStore
}

func NewEnforcer(store RBACStore) *Enforcer {
	return &Enforcer{store: store}
}

func (e *Enforcer) Can(userID string, action Action, resource Resource) bool {
	perms := e.store.GetPermissions(userID)
	for _, p := range perms {
		if p.Action == action && p.Resource == resource {
			return true
		}
	}
	return false
}

func (e *Enforcer) CanOwn(userID string, action Action, resource Resource, ownerID string) bool {
	if userID == ownerID {
		return e.Can(userID, action, resource)
	}
	return false
}
