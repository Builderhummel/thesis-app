package protected_controller

type TableOpenRequests struct {
	ThesisTypeBadge string // BA=bg_primary, MA=...
	ThesisType      string // Type: BA, MA, PA
	Name            string // Name of thesis candiate
	ThesisTitle     string
	RequestDate     string // Unix Time
	Semester        string // Format: SoSe22, WiSe21/22,
	Status          string // requested only xDD
	LinkEmail       string // link to mailto:...
	LinkModify      string
	LinkDelete      string
}

func generateTORTestData() []TableOpenRequests {
	return []TableOpenRequests{
		{
			ThesisTypeBadge: "bg-primary",
			ThesisType:      "BA",
			Name:            "Max Mustermann",
			ThesisTitle:     "Lorem Ipsum",
			RequestDate:     "1630000000",
			Semester:        "SoSe22",
			Status:          "request",
			LinkEmail:       "mailto:abc@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-primary",
			ThesisType:      "BA",
			Name:            "Anna Schmidt",
			ThesisTitle:     "Understanding Data Science",
			RequestDate:     "1622505600",
			Semester:        "WiSe21/22",
			Status:          "request",
			LinkEmail:       "mailto:anna@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-teal",
			ThesisType:      "PA",
			Name:            "John Doe",
			ThesisTitle:     "Quantum Computing Applications",
			RequestDate:     "1635000000",
			Semester:        "SoSe21",
			Status:          "request",
			LinkEmail:       "mailto:john.doe@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-teal",
			ThesisType:      "PA",
			Name:            "Sophie Klein",
			ThesisTitle:     "AI in Healthcare",
			RequestDate:     "1636305600",
			Semester:        "WiSe20/21",
			Status:          "request",
			LinkEmail:       "mailto:sophie.klein@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-primary",
			ThesisType:      "BA",
			Name:            "Lukas Müller",
			ThesisTitle:     "Exploring the Web of Things",
			RequestDate:     "1619702400",
			Semester:        "SoSe20",
			Status:          "request",
			LinkEmail:       "mailto:lukas@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-teal",
			ThesisType:      "PA",
			Name:            "Eva Sommer",
			ThesisTitle:     "The Role of Blockchain in Finance",
			RequestDate:     "1642905600",
			Semester:        "WiSe22/23",
			Status:          "request",
			LinkEmail:       "mailto:eva.sommer@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-primary",
			ThesisType:      "BA",
			Name:            "Michael Bauer",
			ThesisTitle:     "Cloud Computing Security",
			RequestDate:     "1609459200",
			Semester:        "WiSe19/20",
			Status:          "request",
			LinkEmail:       "mailto:michael.bauer@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-primary",
			ThesisType:      "BA",
			Name:            "Julia Becker",
			ThesisTitle:     "Data Privacy in Modern Systems",
			RequestDate:     "1625097600",
			Semester:        "SoSe21",
			Status:          "request",
			LinkEmail:       "mailto:julia.becker@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-purple",
			ThesisType:      "MA",
			Name:            "Felix Wagner",
			ThesisTitle:     "Advances in Computer Vision",
			RequestDate:     "1646092800",
			Semester:        "WiSe22/23",
			Status:          "request",
			LinkEmail:       "mailto:felix.wagner@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
		{
			ThesisTypeBadge: "bg-teal",
			ThesisType:      "PA",
			Name:            "Clara Richter",
			ThesisTitle:     "Automating Data Analysis with AI",
			RequestDate:     "1650403200",
			Semester:        "SoSe22",
			Status:          "request",
			LinkEmail:       "mailto:clara.richter@example.org",
			LinkModify:      "#",
			LinkDelete:      "#",
		},
	}
}
