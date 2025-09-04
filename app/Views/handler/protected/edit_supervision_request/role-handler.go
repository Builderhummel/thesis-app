package view_protected_edit_supervision_request

type ToggleableFields struct {
	Management bool
}

func NewToggleableFields() *ToggleableFields {
	return &ToggleableFields{}
}

func (f *ToggleableFields) ToggleManagementFields(management bool) {
	f.Management = management
}
