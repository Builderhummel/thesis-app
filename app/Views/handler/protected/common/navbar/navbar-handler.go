package view_protected_common_navbar

import "github.com/Builderhummel/thesis-app/app/Constants/roles"

type NavVisibilityState struct {
	Administration bool
}

func NewNavVisibilityState(userRole roles.Role) NavVisibilityState {
	return NavVisibilityState{
		Administration: userRole == roles.RoleAdministrator,
	}
}
