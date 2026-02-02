package protected_controller

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/gin-gonic/gin"
)

type SupervisionRequestForm struct {
	Name          string
	Email         string
	CourseOfStudy string
	ThesisType    string
	ThesisTitle   string
	GPA           string
	ContactDate   string
	Notes         string
}

func RenderAddSupervisionRequestForm(c *gin.Context) {
	c.HTML(http.StatusOK, "protected/add_supervision_request/index.html", gin.H{
		"Navbar": renderNavbar(auth_controller.GetUserRoleFromContext(c)),
	})
}

// TODO: Sanitize input -> why on earth is XSS not working?!?!
func HandleNewSupervisionRequest(c *gin.Context) {
	svrf := SupervisionRequestForm{}
	svrf.Name = c.PostForm("name")
	svrf.Email = c.PostForm("email")
	svrf.CourseOfStudy = c.PostForm("course-of-study")
	svrf.ThesisType = c.PostForm("thesis-type")
	svrf.ThesisTitle = c.PostForm("thesis-title")
	svrf.GPA = c.PostForm("gpa")
	svrf.ContactDate = c.PostForm("contact-date")
	svrf.Notes = c.PostForm("notes")

	// Store values in DB
	tuid, err := db_model.InsertNewThesisRequest(svrf.Name, svrf.Email, svrf.CourseOfStudy, svrf.ThesisType, svrf.ThesisTitle, svrf.GPA, svrf.ContactDate, svrf.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}
	c.Redirect(http.StatusSeeOther, "/view?tuid="+tuid)
}
