package view_protected_edit_supervision_request

type FieldStudentInfo struct {
	Name          string
	Email         string
	CourseOfStudy string
	GPA           string
	Booked        bool
}

func NewFieldStudentInfo() FieldStudentInfo {
	return FieldStudentInfo{}
}

func (f *FieldStudentInfo) SetInfo(Name string, Email string, CourseOfStudy string, GPA string, Booked bool) {
	f.Name = Name
	f.Email = Email
	f.CourseOfStudy = CourseOfStudy
	f.GPA = GPA
	f.Booked = Booked
}
