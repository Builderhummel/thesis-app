package view_protected_view_supervision_request

import (
	"fmt"
	"time"
)

type FieldThesisInfo struct {
	Tuid         string
	ThesisType   string
	ThesisTitle  string
	ThesisStatus string
	FinalGrade   string
	RequestDate  string
	ContactDate  string
	Deadline     string
	SubmitDate   string
	Supervisors  []string
	Examiners    []string
	Notes        string
}

func NewFieldThesisInfo() FieldThesisInfo {
	return FieldThesisInfo{}
}

func (f *FieldThesisInfo) SetInfo(Tuid, ThesisType, ThesisTitle, ThesisStatus string, FinalGrade float64, RequestDate, ContactDate, Deadline, SubmitDate time.Time, Supervisors []string, Examiners []string, Notes string) {
	f.Tuid = Tuid
	f.ThesisType = ThesisType
	f.ThesisTitle = ThesisTitle
	f.ThesisStatus = ThesisStatus
	f.FinalGrade = overwriteFinalGrade(FinalGrade)
	f.RequestDate = overwriteZeroDate(RequestDate, "n/a")
	f.ContactDate = overwriteZeroDate(ContactDate, "Not contacted yet")
	f.Deadline = overwriteZeroDate(Deadline, "No deadline set")
	f.SubmitDate = overwriteZeroDate(SubmitDate, "Not submitted yet")
	f.Supervisors = overwriteEmptySlice(Supervisors, "No supervisor assigned")
	f.Examiners = overwriteEmptySlice(Examiners, "No examiner assigned")
	f.Notes = Notes
}

func overwriteZeroDate(date time.Time, overwrite string) string {
	if date.IsZero() {
		return overwrite
	}
	return date.Format("2006-01-02")
}

func overwriteFinalGrade(grade float64) string {
	if grade == -1 {
		return "Not graded yet"
	}
	return fmt.Sprintf("%.2f", grade)
}

func overwriteEmptySlice(slice []string, overwrite string) []string {
	if len(slice) == 0 {
		return []string{overwrite}
	}
	return slice
}
