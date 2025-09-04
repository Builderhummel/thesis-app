package protected_controller

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_view_supervision_request "github.com/Builderhummel/thesis-app/app/Views/handler/protected/view_supervision_request"
	"github.com/gin-gonic/gin"
)

func RenderViewSupervisionRequestForm(c *gin.Context) {
	referer := c.Request.Referer()

	tuid := html.EscapeString(c.Query("tuid"))
	if tuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tuid provided"})
		return
	}

	if tuid_num, err := strconv.Atoi(tuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tuid"})
		return
	} else if tuid_num < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tuid"})
		return
	}

	tfd, err := db_model.GetDataFullSupervision(tuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}

	studInf := view_protected_view_supervision_request.NewFieldStudentInfo()
	studInf.SetInfo(tfd.Name, tfd.Email, tfd.StudyProgram, fmt.Sprintf("%.2f", tfd.GPA))

	thesisInf := view_protected_view_supervision_request.NewFieldThesisInfo()
	thesisInf.SetInfo(tuid, tfd.ThesisType, tfd.ThesisTitle, tfd.ThesisStatus, tfd.FinalGrade, tfd.RequestDate, tfd.ResponseDate, tfd.RegisteredDate, tfd.Deadline, tfd.SubmitDate, getSupervisorsSliceFromPersonalData(tfd.Supervisors), getSupervisorsSliceFromPersonalData(tfd.Examiners), tfd.Notes)

	c.HTML(http.StatusOK, "protected/view_supervision_request/index.html", gin.H{
		"Navbar":    renderNavbar(auth_controller.GetUserRoleFromContext(c)),
		"Referer":   referer,
		"StudInf":   studInf,
		"ThesisInf": thesisInf,
	})
}
