package protected_controller

import (
	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_all_supervisions "github.com/Builderhummel/thesis-app/app/Views/handler/protected/all_supervisions"
	"github.com/gin-gonic/gin"
)

func RenderAllSupervisions(c *gin.Context) {
	tabAllSupervisions, err := fillTableAllSupervisions()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "protected/all_supervisions/index.html", gin.H{
		"TabAllSupV": tabAllSupervisions,
		"Navbar":     renderNavbar(auth_controller.GetUserRoleFromContext(c)),
	})
}

func fillTableAllSupervisions() (view_protected_all_supervisions.TableAllSupervisions, error) {
	tab_data, err := db_model.GetDataTableAllSupervisions()
	if err != nil {
		return nil, err
	}
	tasv := view_protected_all_supervisions.NewTableMySupervisions()
	for _, row_data := range tab_data {
		tasv.AddRow(row_data["thesisType"], row_data["name"], row_data["thesisTitle"], row_data["deadline"], row_data["supervisor"], row_data["semester"], row_data["status"], "mailto:"+row_data["email"], "/view?tuid="+row_data["tuid"], "/delete?tuid="+row_data["tuid"])
	}
	return tasv, nil
}
