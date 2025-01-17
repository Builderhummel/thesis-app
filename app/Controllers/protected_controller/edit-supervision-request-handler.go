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
		tuid,
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
		"Navbar":    renderNavbar(),
		"StudInf":   studInf,
		"ThesisInf": thesisInf,
	})
}

func HandleEditSupervisionRequest(c *gin.Context) {
	tfd := db_model.ThesisFullData{}
	supervisors := []db_model.PersonalData{}
	examiners := []db_model.PersonalData{}

	tuid := c.Query("tuid")
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

	gpa, err := processGradeString(c.PostForm("gpa"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with GPA"})
		return
	}

	finalGrade, err := processGradeString(c.PostForm("final-grade"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with final grade"})
		return
	}

	requestDate, err := parseDateStringToGoDate(c.PostForm("request-date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with request date"})
		return
	}

	contactDate, err := parseDateStringToGoDate(c.PostForm("contact-date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with contact date"})
		return
	}

	deadline, err := parseDateStringToGoDate(c.PostForm("deadline"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with deadline"})
		return
	}

	submitDate, err := parseDateStringToGoDate(c.PostForm("submit-date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with submit date"})
		return
	}

	tfd.TUID = tuid
	tfd.Name = c.PostForm("name")
	tfd.Email = c.PostForm("email")
	tfd.StudyProgram = c.PostForm("study-program")
	tfd.GPA = gpa
	tfd.Booked = c.PostForm("thesis-booked") == "true"
	tfd.ThesisType = c.PostForm("thesis-type")
	tfd.ThesisTitle = c.PostForm("thesis-title")
	tfd.ThesisStatus = c.PostForm("thesis-status")
	tfd.Semester = concatSemesterInfo(c.PostForm("thesis-semester"), c.PostForm("thesis-semester-year")) //to handle thesis-semester, thesis-semester-year
	tfd.FinalGrade = finalGrade
	tfd.RequestDate = requestDate
	tfd.ContactDate = contactDate
	tfd.Deadline = deadline
	tfd.SubmitDate = submitDate
	tfd.Notes = c.PostForm("notes")

	supervisors_puid := c.PostFormArray("supervisors[]")
	examiners_puid := c.PostFormArray("examiners[]")

	for _, puid := range supervisors_puid {
		supervisors = append(supervisors, db_model.PersonalData{PDUid: puid})
	}
	for _, puid := range examiners_puid {
		examiners = append(examiners, db_model.PersonalData{PDUid: puid})
	}

	tfd.Supervisors = supervisors
	tfd.Examiners = examiners

	err = db_model.UpdateThesisInfo(tfd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}

	returnUrl := fmt.Sprintf("/view?tuid=%s", tuid)
	c.Redirect(http.StatusSeeOther, returnUrl) // TODO: Add success message (flash alert?)
}
