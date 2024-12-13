package protected_controller

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_edit_supervision_request "github.com/Builderhummel/thesis-app/app/Views/handler/protected/edit_supervision_request"
	"github.com/gin-gonic/gin"
)

func RenderEditSupervisionRequestForm(c *gin.Context) {
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

	allSupervisors, err := db_model.GetAllSupervisors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}
	annotatedSupervisors := annotateSelectedPDs(
		convertSliceDataPDtoViewEditPD(allSupervisors),
		convertSliceDataPDtoViewEditPD(tfd.Supervisors))

	allExaminers, err := db_model.GetAllExaminers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}
	annotatedExaminers := annotateSelectedPDs(
		convertSliceDataPDtoViewEditPD(allExaminers),
		convertSliceDataPDtoViewEditPD(tfd.Examiners))

	studInf := view_protected_edit_supervision_request.NewFieldStudentInfo()
	studInf.SetInfo(tfd.Name, tfd.Email, tfd.StudyProgram, fmt.Sprintf("%.2f", tfd.GPA), tfd.Booked)

	thesisInf := view_protected_edit_supervision_request.NewFieldThesisInfo()
	thesisInf.SetInfo(
		tfd.ThesisType,
		tfd.ThesisTitle,
		tfd.ThesisStatus, tfd.Semester,
		tfd.FinalGrade,
		tfd.RequestDate,
		tfd.ContactDate,
		tfd.Deadline, tfd.SubmitDate,
		annotatedSupervisors,
		annotatedExaminers,
		tfd.Notes)

	c.HTML(http.StatusOK, "protected/edit_supervision_request/index.html", gin.H{
		"StudInf":   studInf,
		"ThesisInf": thesisInf,
	})
}

func HandleEditSupervisionRequest(c *gin.Context) {
	fmt.Printf("a%+v\n", c.PostFormArray("supervisors[]"))

	c.JSON(http.StatusOK, gin.H{"message": "Success!"})
}
