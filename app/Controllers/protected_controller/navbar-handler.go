package protected_controller

import (
	"github.com/Builderhummel/thesis-app/app/Constants/roles"
	view_protected_common_navbar "github.com/Builderhummel/thesis-app/app/Views/handler/protected/common/navbar"
)

type Navbar struct {
	view_protected_common_navbar.NavVisibilityState
}

func renderNavbar(userRole roles.Role) Navbar {
	navbarVisState := view_protected_common_navbar.NewNavVisibilityState(userRole)
	return Navbar{
		navbarVisState}
}
