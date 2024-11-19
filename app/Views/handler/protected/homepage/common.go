package view_protected_homepage

func setThesisTypeBadge(thesisType string) string {
	switch thesisType {
	case "BA":
		return "bg-primary"
	case "MA":
		return "bg-purple"
	case "PA":
		return "bg-teal"
	default:
		return ""
	}
}
