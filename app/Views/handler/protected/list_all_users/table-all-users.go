package listallusers

type TableRowAllUsers struct {
	Name         string
	Email        string
	LoginHandle  string
	IsActive     bool
	IsSupervisor bool
	IsExaminer   bool
}

type TableAllUsers []TableRowAllUsers

func NewTableAllUsers() TableAllUsers {
	return TableAllUsers{}
}

func (t *TableAllUsers) AddRow(Name, Email, LoginHandle string, IsActive, IsSupervisor, IsExaminer bool) {
	*t = append(*t, NewTableRowAllUsers(Name, Email, LoginHandle, IsActive, IsSupervisor, IsExaminer))
}

func NewTableRowAllUsers(Name, Email, LoginHandle string, IsActive, IsSupervisor, IsExaminer bool) TableRowAllUsers {
	return TableRowAllUsers{
		Name:         Name,
		Email:        Email,
		LoginHandle:  LoginHandle,
		IsActive:     IsActive,
		IsSupervisor: IsSupervisor,
		IsExaminer:   IsExaminer,
	}
}
