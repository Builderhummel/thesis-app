package view_protected_all_supervisions

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

// request, contacted, registered, working, completed, dropped
// TODO: Change colors
func setStatusBadge(statusBadge string) string {
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
	case "reject":
		return "bg-dark"
	case "dropped":
		return "bg-dark"
	default:
		return ""
	}
}
