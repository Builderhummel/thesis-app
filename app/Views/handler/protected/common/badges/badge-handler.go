package view_protected_common_badges

func SetThesisTypeBadge(thesisType string) string {
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

// request, contacted, registered, working, completed, dropped
// TODO: Change colors
func SetStatusBadge(statusBadge string) string {
	switch statusBadge {
	case "request":
		return "bg-danger"
	case "contacted":
		return "bg-warning"
	case "registered":
		return "bg-info"
	case "working":
		return "bg-success"
	case "completed":
		return "bg-secondary"
	case "dropped":
		return "bg-dark"
	case "reject":
		return "bg-dark"
	default:
		return ""
	}
}
