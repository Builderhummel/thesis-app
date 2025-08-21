package view_protected_edit_user

type FieldUser struct {
	PDUID        string
	Name         string
	Email        string
	LoginHandle  string
	Role         string
	IsActive     bool
	IsSupervisor bool
	IsExaminer   bool
}

func NewFieldUser() FieldUser {
	return FieldUser{}
}

func (f *FieldUser) SetUser(pduid string, name string, email string, loginHandle string, role string, isActive bool, isSupervisor bool, isExaminer bool) {
	f.PDUID = pduid
	f.Name = name
	f.Email = email
	f.LoginHandle = loginHandle
	f.Role = role
	f.IsActive = isActive
	f.IsSupervisor = isSupervisor
	f.IsExaminer = isExaminer
}
