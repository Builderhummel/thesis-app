package view_protected_listallusers

type TableRowAllUsers struct {
	PDUID        string
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

func (t *TableAllUsers) AddRow(PDUID, Name, Email, LoginHandle string, IsActive, IsSupervisor, IsExaminer bool) {
	*t = append(*t, NewTableRowAllUsers(PDUID, Name, Email, LoginHandle, IsActive, IsSupervisor, IsExaminer))
}

func NewTableRowAllUsers(PDUID, Name, Email, LoginHandle string, IsActive, IsSupervisor, IsExaminer bool) TableRowAllUsers {
	return TableRowAllUsers{
		PDUID:        PDUID,
		Name:         Name,
		Email:        Email,
		LoginHandle:  LoginHandle,
		IsActive:     IsActive,
		IsSupervisor: IsSupervisor,
		IsExaminer:   IsExaminer,
	}
}
