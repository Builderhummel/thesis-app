package protected_controller

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_my_supervisions "github.com/Builderhummel/thesis-app/app/Views/handler/protected/my_supervisions"
	"github.com/gin-gonic/gin"
)

func RenderMySupervisions(c *gin.Context) {
	//Get User handle
	user_id, err := auth_controller.ExtractTokenUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract user handle"}) //TODO: Proper error handling
		return
	}

	tabMySupV, err := fillTableMySupervisions(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}
	c.HTML(http.StatusOK, "protected/my_supervisions/index.html", gin.H{
		"Navbar":    renderNavbar(),
		"TabMySupV": tabMySupV,
	})
}

func fillTableMySupervisions(user_id string) (view_protected_my_supervisions.TableMySupervisions, error) {
	tab_data, err := db_model.GetDataTableMySupervisions(user_id)
	if err != nil {
		return nil, err
	}
	tmsv := view_protected_my_supervisions.NewTableMySupervisions()
	for _, row_data := range tab_data {
		tmsv.AddRow(row_data["thesistype"], row_data["name"], row_data["thesistitle"], row_data["deadline"], row_data["supervisor"], row_data["semester"], row_data["thesisstatus"], "mailto:"+row_data["email"], "/view?tuid="+row_data["tuid"], "#")
	}
	return tmsv, nil
}
