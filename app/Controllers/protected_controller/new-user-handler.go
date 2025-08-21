package protected_controller

import (
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/gin-gonic/gin"
)

func RenderNewUser(c *gin.Context) {
	c.HTML(200, "protected/new_user/index.html", gin.H{
		"Navbar": renderNavbar(),
	})
}

func HandlePostNewUser(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	handle := c.PostForm("handle")
	role := c.PostForm("role")
	active := c.PostForm("active") == "on"
	isSupervisor := c.PostForm("isSupervisor") == "on"
	isExaminer := c.PostForm("isExaminer") == "on"

	err := db_model.InsertNewUser(name, email, handle, role, active, isSupervisor, isExaminer)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error inserting new user: " + err.Error()})
		return
	}

	c.Redirect(302, "/users")
}
