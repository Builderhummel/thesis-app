package view_protected_view_supervision_request

type FieldStudentInfo struct {
	Name          string
	Email         string
	CourseOfStudy string
	GPA           string
}

func NewFieldStudentInfo() FieldStudentInfo {
	return FieldStudentInfo{}
}

func (f *FieldStudentInfo) SetInfo(Name string, Email string, CourseOfStudy string, GPA string) {
	f.Name = Name
	f.Email = Email
	f.CourseOfStudy = CourseOfStudy
	f.GPA = GPA
}
