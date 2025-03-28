package view_protected_open_requests

import (
	view_protected_common_badges "github.com/Builderhummel/thesis-app/app/Views/handler/protected/common/badges"
)

type TableOpenRequest []TableRowOpenRequest

type TableRowOpenRequest struct {
	ThesisTypeBadge string // BA=bg_primary, MA=...
	ThesisType      string // Type: BA, MA, PA
	Name            string // Name of thesis candiate
	CourseOfStudy   string
	GPA             string // e.g. 4.0
	RequestDate     string // Unix Time
	Status          string // requested only xDD
	LinkEmail       string // link to mailto:...
	LinkModify      string
	LinkDelete      string
}

func NewTableOpenRequests() TableOpenRequest {
	return TableOpenRequest{}
}

func (t *TableOpenRequest) AddRow(ThesisType, Name, ThesisTitle, RequestDate, Semester, Status, LinkEmail, LinkModify, LinkDelete string) {
	*t = append(*t, NewTableRowOpenRequest(ThesisType, Name, ThesisTitle, RequestDate, Semester, Status, LinkEmail, LinkModify, LinkDelete))
}

func NewTableRowOpenRequest(ThesisType, Name, CourseOfStudy, GPA, RequestDate, Status, LinkEmail, LinkModify, LinkDelete string) TableRowOpenRequest {
	return TableRowOpenRequest{
		ThesisTypeBadge: view_protected_common_badges.SetThesisTypeBadge(ThesisType),
		ThesisType:      ThesisType,
		Name:            Name,
		CourseOfStudy:   CourseOfStudy,
		GPA:             GPA,
		RequestDate:     RequestDate,
		Status:          Status,
		LinkEmail:       LinkEmail,
		LinkModify:      LinkModify,
		LinkDelete:      LinkDelete,
	}
}
