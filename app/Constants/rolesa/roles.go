package roles

type Role string

const (
	RoleDefault       Role = "default" // No rights, every user will get this role if not specified
	RoleResearcher    Role = "researcher"
	RoleManagement    Role = "management"
	RoleAdministrator Role = "administrator"
)
