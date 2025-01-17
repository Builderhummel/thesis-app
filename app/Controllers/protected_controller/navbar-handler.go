package protected_controller

import view_protected_common_navbar "github.com/Builderhummel/thesis-app/app/Views/handler/protected/common/navbar"

type Navbar struct {
	view_protected_common_navbar.NavVisibilityState
}

func renderNavbar() Navbar {
	navbarVisState := view_protected_common_navbar.NewNavVisibilityState()
	navbarVisState.Init()
	return Navbar{
		navbarVisState}
}
