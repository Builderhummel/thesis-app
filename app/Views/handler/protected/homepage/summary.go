package view_protected_homepage

type Summary struct {
	Requested string
	Contacted string
	Working   string
}

func NewSummary() Summary {
	return Summary{}
}

func (s *Summary) SetSummary(Requested, Contacted, Working string) {
	s.Requested = Requested
	s.Contacted = Contacted
	s.Working = Working
}
