package view_protected_all_supervisions

type TableAllSupervisions []TableRowAllSupervisions

type TableRowAllSupervisions struct {
	ThesisTypeBadge string // BA=bg_primary, MA=...
	ThesisType      string // Type: BA, MA, PA
	Name            string // Name of thesis candiate
	ThesisTitle     string
	DeadlineDate    string // Unix Time
	Supervisor      string // List of Supervisors
	Semester        string // Format: SoSe22, WiSe21/22,
	Status          string // everything but request
	StatusBadge     string
	LinkEmail       string // link to mailto:...
	LinkModify      string
	LinkDelete      string
}

func NewTableMySupervisions() TableAllSupervisions {
	return TableAllSupervisions{}
}

func (t *TableAllSupervisions) AddRow(ThesisType, Name, ThesisTitle, DeadlineDate, Supervisor, Semester, Status, LinkEmail, LinkModify, LinkDelete string) {
	*t = append(*t, NewTableRowMySupervisions(ThesisType, Name, ThesisTitle, DeadlineDate, Supervisor, Semester, Status, LinkEmail, LinkModify, LinkDelete))
}

func NewTableRowMySupervisions(ThesisType, Name, ThesisTitle, DeadlineDate, Supervisor, Semester, Status, LinkEmail, LinkModify, LinkDelete string) TableRowAllSupervisions {
	return TableRowAllSupervisions{
		ThesisTypeBadge: setThesisTypeBadge(ThesisType),
		ThesisType:      ThesisType,
		Name:            Name,
		ThesisTitle:     ThesisTitle,
		DeadlineDate:    DeadlineDate,
		Supervisor:      Supervisor,
		Semester:        Semester,
		Status:          Status,
		StatusBadge:     setStatusBadge(Status),
		LinkEmail:       LinkEmail,
		LinkModify:      LinkModify,
		LinkDelete:      LinkDelete,
	}
}
