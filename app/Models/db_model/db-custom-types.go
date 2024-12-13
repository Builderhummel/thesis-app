package db_model

import "time"

type ThesisFullData struct {
	TUID         string
	Name         string
	Email        string
	StudyProgram string
	Booked       bool
	GPA          float64
	ThesisType   string
	ThesisTitle  string
	ThesisStatus string
	Semester     string
	FinalGrade   float64
	RequestDate  time.Time
	ContactDate  time.Time
	Deadline     time.Time
	SubmitDate   time.Time
	Supervisors  []PersonalData
	Examiners    []PersonalData
	Notes        string
}

type PersonalData struct {
	PDUid string
	Name  string
	Email string
}
