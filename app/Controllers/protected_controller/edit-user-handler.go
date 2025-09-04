package protected_controller

import (
	"net/http"
	"strconv"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_admin_edit_user "github.com/Builderhummel/thesis-app/app/Views/handler/admin/edit_user"
	"github.com/gin-gonic/gin"
)

func RenderEditUser(c *gin.Context) {
	puid := c.Query("puid")
	if puid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No puid provided"})
		return
	}

	if puid_num, err := strconv.Atoi(puid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid puid"})
		return
	} else if puid_num < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid puid"})
		return
	}

	userData, err := db_model.GetUserByPUID(puid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user data" + err.Error()})
		return
	}

	fUser := view_admin_edit_user.NewFieldUser()
	fUser.SetUser(userData.PDUid, userData.Name, userData.Email, userData.Handle, string(userData.Role), userData.IsActive, userData.IsSupervisor, userData.IsExaminer)

	c.HTML(200, "admin/edit_user/index.html", gin.H{
		"Navbar": renderNavbar(auth_controller.GetUserRoleFromContext(c)),
		"User":   fUser,
	})
}

func HandlePostEditUser(c *gin.Context) {
	puid := c.Query("puid")

	name := c.PostForm("name")
	email := c.PostForm("email")
	handle := c.PostForm("handle")
	role := c.PostForm("role")
	active := c.PostForm("active") == "on"
	isSupervisor := c.PostForm("isSupervisor") == "on"
	isExaminer := c.PostForm("isExaminer") == "on"

	if puid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No puid provided"})
		return
	}
	if puid_num, err := strconv.Atoi(puid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid puid"})
		return
	} else if puid_num < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid puid"})
		return
	}

	err := db_model.UpdateFullUser(puid, name, email, handle, role, active, isSupervisor, isExaminer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user" + err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/users")
}
