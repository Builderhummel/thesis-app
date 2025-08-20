package protected_controller

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/gin-gonic/gin"
)

// Gin handler
func HandleExport(c *gin.Context) {
	format := c.Query("format")
	switch format {
	case "csv":
		data, err := db_model.GetAllDataFullSupervison()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		writeThesisCSV(c, data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format not supported/specified"})
	}
}

// Utility: Converts []PersonalData to "Name1|Name2|Name3"
func joinNames(people []db_model.PersonalData) string {
	names := make([]string, len(people))
	for i, p := range people {
		names[i] = p.Name
	}
	return strings.Join(names, "|")
}

// Utility: Converts ThesisFullData to []string for CSV
func thesisToCSVRow(t *db_model.ThesisFullData) []string {
	// Use strconv for numbers, time.Time.Format for dates, etc.
	layout := "2006-01-02"
	return []string{
		t.TUID,
		t.Name,
		t.Email,
		t.StudyProgram,
		strconv.FormatBool(t.Booked),
		strconv.FormatFloat(t.GPA, 'f', 2, 64),
		t.ThesisType,
		t.ThesisTitle,
		t.ThesisStatus,
		t.Semester,
		strconv.FormatFloat(t.FinalGrade, 'f', 2, 64),
		t.RequestDate.Format(layout),
		t.ResponseDate.Format(layout),
		t.RegisteredDate.Format(layout),
		t.Deadline.Format(layout),
		t.SubmitDate.Format(layout),
		joinNames(t.Supervisors),
		joinNames(t.Examiners),
		t.Notes,
	}
}

// Utility: CSV header
func csvHeader() []string {
	return []string{
		"TUID", "Name", "Email", "StudyProgram", "Booked", "GPA",
		"ThesisType", "ThesisTitle", "ThesisStatus", "Semester", "FinalGrade",
		"RequestDate", "ResponseDate", "RegisteredDate", "Deadline", "SubmitDate",
		"Supervisors", "Examiners", "Notes",
	}
}

// Handles writing CSV to gin.Context
func writeThesisCSV(c *gin.Context, data []*db_model.ThesisFullData) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=thesisinfo_full_export_"+time.Now().Format("2006_01_02_1504")+".csv")
	w := csv.NewWriter(c.Writer)
	_ = w.Write(csvHeader())

	for _, t := range data {
		_ = w.Write(thesisToCSVRow(t))
	}
	w.Flush()
}
