package view_protected_delete_supervision

type ThesisInfo struct {
	Tuid        string
	ThesisTitle string
}

type ViewDelete struct {
	Name       string
	Email      string
	ThesisInfo ThesisInfo
}

func FillDeleteSupervision(
	name string,
	email string,
	tuid string,
	thesisTitle string,
) ViewDelete {
	return ViewDelete{
		Name:  name,
		Email: email,
		ThesisInfo: ThesisInfo{
			Tuid:        tuid,
			ThesisTitle: thesisTitle,
		},
	}
}
