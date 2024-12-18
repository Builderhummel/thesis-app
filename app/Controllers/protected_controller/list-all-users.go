package protected_controller

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	listallusers "github.com/Builderhummel/thesis-app/app/Views/handler/protected/list_all_users"
	"github.com/gin-gonic/gin"
)

func RenderListAllUsers(c *gin.Context) {

	allUsers, err := db_model.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	tabAllUsrs := listallusers.NewTableAllUsers()
	for _, user := range allUsers {
		tabAllUsrs.AddRow(user.Name, user.Email, user.Handle, user.IsActive, user.IsSupervisor, user.IsExaminer)
	}

	c.HTML(http.StatusOK, "protected/list_all_users/index.html", gin.H{
		"Users": tabAllUsrs,
	})
}
