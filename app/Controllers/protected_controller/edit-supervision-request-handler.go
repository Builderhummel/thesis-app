package protected_controller

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/Builderhummel/thesis-app/app/Constants/roles"
	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
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

	userrole := auth_controller.GetUserRoleFromContext(c)

	toggleableFields := view_protected_edit_supervision_request.NewToggleableFields()
	toggleableFields.ToggleManagementFields(auth_controller.MinUserGroup(userrole, roles.RoleManagement))

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
		tfd.GitlabRepo,
		tfd.RequestDate,
		tfd.ResponseDate,
		tfd.RegisteredDate,
		tfd.Deadline, tfd.SubmitDate,
		annotatedSupervisors,
		annotatedExaminers,
		tfd.Notes)

	c.HTML(http.StatusOK, "protected/edit_supervision_request/index.html", gin.H{
		"Navbar":           renderNavbar(userrole),
		"ToggleableFields": toggleableFields,
		"StudInf":          studInf,
		"ThesisInf":        thesisInf,
	})
}

func HandleEditSupervisionRequest(c *gin.Context) {
	userRole := auth_controller.GetUserRoleFromContext(c)

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

	// Get original thesis data
	tfd, err := db_model.GetDataFullSupervision(tuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
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

	responseDate, err := parseDateStringToGoDate(c.PostForm("contact-date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with contact date"})
		return
	}

	registeredDate, err := parseDateStringToGoDate(c.PostForm("registered-date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error with registered date"})
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

	// Update the data according to user role permissions
	// Min role: Default
	if auth_controller.MinUserGroup(userRole, roles.RoleDefault) {
		// Default has to rights atm
	}

	// Min role: Researcher
	if auth_controller.MinUserGroup(userRole, roles.RoleResearcher) {
		tfd.TUID = tuid
		tfd.Name = c.PostForm("name")
		tfd.Email = c.PostForm("email")
		tfd.StudyProgram = c.PostForm("study-program")
		tfd.GPA = gpa
		tfd.ThesisType = c.PostForm("thesis-type")
		tfd.ThesisTitle = c.PostForm("thesis-title")
		tfd.ThesisStatus = c.PostForm("thesis-status")
		tfd.FinalGrade = finalGrade
		tfd.GitlabRepo = c.PostForm("gitlab-repo")
		tfd.RequestDate = requestDate
		tfd.ResponseDate = responseDate
		tfd.RegisteredDate = registeredDate
		tfd.Deadline = deadline
		tfd.SubmitDate = submitDate
		tfd.Notes = c.PostForm("notes")

		//Supervision Fields
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
	}

	// Min role: Management
	if auth_controller.MinUserGroup(userRole, roles.RoleManagement) {
		tfd.Semester = concatSemesterInfo(c.PostForm("thesis-semester"), c.PostForm("thesis-semester-year")) //to handle thesis-semester, thesis-semester-year
		tfd.Booked = c.PostForm("thesis-booked") == "true"
	}

	// Min role: Administrator
	if auth_controller.MinUserGroup(userRole, roles.RoleAdministrator) {
		//TODO: Only things an administrator can change
	}

	err = db_model.UpdateThesisInfo(tfd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}

	returnUrl := fmt.Sprintf("/view?tuid=%s", tuid)
	c.Redirect(http.StatusSeeOther, returnUrl) // TODO: Add success message (flash alert?)
}

// POST function
func HandleEditAssignToMe(c *gin.Context) {
	user_id, err := auth_controller.ExtractTokenUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract user handle"}) //TODO: Proper error handling
		return
	}

	userRole := auth_controller.GetUserRoleFromContext(c)
	if !auth_controller.MinUserGroup(userRole, roles.RoleResearcher) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	tuid := html.EscapeString(c.PostForm("tuid"))
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

	userData, err := db_model.GetUserByLoginHandle(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user data"})
		return
	}

	err = db_model.AddThesisSupervisor(tuid, userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not assign supervisor"})
		return
	}

	returnUrl := fmt.Sprintf("/view?tuid=%s", tuid)
	c.Redirect(http.StatusSeeOther, returnUrl)
}
