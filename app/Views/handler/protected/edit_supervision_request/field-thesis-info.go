package view_protected_edit_supervision_request

import (
	"fmt"
	"regexp"
	"time"
)

type FieldThesisInfo struct {
	Tuid              string
	ThesisType        string
	ThesisTitle       string
	ThesisStatus      string
	SemesterSelection string //1: WiSe, 2: SoSe
	SemesterYear      string
	FinalGrade        string
	RequestDate       string
	ContactDate       string
	Deadline          string
	SubmitDate        string
	Supervisors       []PersonalData
	Examiners         []PersonalData
	Notes             string
}

type PersonalData struct {
	PDUid    string
	Name     string
	Selected bool
}

func NewFieldThesisInfo() FieldThesisInfo {
	return FieldThesisInfo{}
}

func NewPersonalData() PersonalData {
	return PersonalData{}
}

func (f *FieldThesisInfo) SetInfo(Tuid, ThesisType, ThesisTitle, ThesisStatus, Semester string, FinalGrade float64, RequestDate, ContactDate, Deadline, SubmitDate time.Time, Supervisors []PersonalData, Examiners []PersonalData, Notes string) {
	f.Tuid = Tuid
	f.ThesisType = ThesisType
	f.ThesisTitle = ThesisTitle
	f.ThesisStatus = ThesisStatus
	f.SemesterSelection = getSemesterSeason(Semester)
	f.SemesterYear = getSemesterYear(Semester)
	f.FinalGrade = overwriteFinalGrade(FinalGrade)
	f.RequestDate = overwriteZeroDate(RequestDate, "")
	f.ContactDate = overwriteZeroDate(ContactDate, "")
	f.Deadline = overwriteZeroDate(Deadline, "")
	f.SubmitDate = overwriteZeroDate(SubmitDate, "")
	f.Supervisors = overwriteEmptyPD(Supervisors)
	f.Examiners = overwriteEmptyPD(Examiners)
	f.Notes = Notes
}

func overwriteZeroDate(date time.Time, overwrite string) string {
	if date.IsZero() {
		return overwrite
	}
	return date.Format("2006-01-02")
}

// TODO: Make sure that its stored as -1 again if its empty
func overwriteFinalGrade(grade float64) string {
	if grade == -1 {
		return ""
	}
	return fmt.Sprintf("%.2f", grade)
}

func overwriteEmptyPD(slice []PersonalData) []PersonalData {
	if len(slice) == 0 {
		return []PersonalData{{Name: "", PDUid: ""}}
	}
	return slice
}

func splitSemesterInfo(semester string) (string, string, error) {
	re := regexp.MustCompile(`^(WiSe|SoSe)(\d{2}/\d{2})$`)

	matches := re.FindStringSubmatch(semester)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid semester format")
	}

	return matches[1], matches[2], nil
}

func getSemesterSeason(semester string) string {
	season, _, err := splitSemesterInfo(semester)
	if err != nil {
		return ""
	}

	return season
}

func getSemesterYear(semester string) string {
	_, year, err := splitSemesterInfo(semester)
	if err != nil {
		return ""
	}

	return year
}
