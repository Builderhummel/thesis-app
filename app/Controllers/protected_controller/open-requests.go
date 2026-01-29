package protected_controller

import (
	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_open_requests "github.com/Builderhummel/thesis-app/app/Views/handler/protected/open_requests"
	"github.com/gin-gonic/gin"
)

func RenderOpenRequests(c *gin.Context) {
	tabOpReq, err := fillTableOpenRequests()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "protected/open_requests/index.html", gin.H{
		"TabOpReq": tabOpReq,
		"Navbar":   renderNavbar(auth_controller.GetUserRoleFromContext(c)),
	})
}

func fillTableOpenRequests() (view_protected_open_requests.TableOpenRequest, error) {
	tor := view_protected_open_requests.NewTableOpenRequests()
	tor_data, err := db_model.GetDataThesisTableOpenRequests()
	if err != nil {
		return nil, err
	}
	for _, row_data := range tor_data {
		tor.AddRow(row_data["thesisType"], row_data["name"], row_data["courseOfStudy"], row_data["gpa"], row_data["requestDate"], row_data["status"], row_data["email"], row_data["tuid"])
	}
	return tor, nil
}
