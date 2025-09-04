package roles

type Role string

const (
	RoleDefault       Role = "default" // No rights, every user will get this role if not specified
	RoleResearcher    Role = "researcher"
	RoleManagement    Role = "management"
	RoleAdministrator Role = "administrator"
)

var RoleHierarchy = map[Role]int{
	RoleDefault:       0,
	RoleResearcher:    1,
	RoleManagement:    2,
	RoleAdministrator: 3,
}
