package protected_controller

import (
	"net/http"
	"strconv"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_delete_supervision "github.com/Builderhummel/thesis-app/app/Views/handler/protected/delete_supervision"
	"github.com/gin-gonic/gin"
)

func RenderDeleteSupervisionRequestForm(c *gin.Context) {
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

	thesisData, err := db_model.GetDataFullSupervision(tuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting thesis data: " + err.Error()})
		return
	}

	if thesisData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thesis request not found"})
		return
	}

	vdeleteData := view_protected_delete_supervision.FillDeleteSupervision(
		thesisData.Name,
		thesisData.Email,
		thesisData.TUID,
		thesisData.ThesisTitle,
	)

	c.HTML(200, "protected/view_delete_supervision_request/index.html", gin.H{
		"Navbar":   renderNavbar(auth_controller.GetUserRoleFromContext(c)),
		"VDelData": vdeleteData,
	})
}

func HandleDeleteThesisRequest(c *gin.Context) {
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

	booked, err := db_model.CheckIfThesisIsBooked(tuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking thesis booking status: " + err.Error()})
		return
	}

	if booked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thesis request is booked and cannot be deleted"})
		return
	}

	if err := db_model.DeleteThesisRequest(tuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete thesis request: " + err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/all")
}
