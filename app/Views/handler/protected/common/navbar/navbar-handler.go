package view_protected_common_navbar

type NavVisibilityState struct {
	Administration bool
}

func NewNavVisibilityState() NavVisibilityState {
	return NavVisibilityState{}
}

func (s *NavVisibilityState) SetStates(Administration bool) {
	s.Administration = Administration
}

func (s *NavVisibilityState) Init() error {
	s.Administration = true
	return nil
}
