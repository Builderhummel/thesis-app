package protected_controller

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_homepage "github.com/Builderhummel/thesis-app/app/Views/handler/protected/homepage"
	"github.com/gin-gonic/gin"
)

func RenderHomepage(c *gin.Context) {
	//Get User handle
	user_id, err := auth_controller.ExtractTokenUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract user handle"}) //TODO: Proper error handling
		return
	}
	//Table Open Requests
	tabOpReq, err := fillTableOpenRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}

	tabMySupV, err := fillTableMySupervisions(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
		return
	}
	c.HTML(http.StatusOK, "protected/homepage/index.html", gin.H{
		"TabOpReq":  tabOpReq,
		"TabMySupV": tabMySupV,
	})
}

func fillTableOpenRequests() (view_protected_homepage.TableOpenRequest, error) {
	tor := view_protected_homepage.NewTableOpenRequests()
	tor_data, err := db_model.GetDataThesisTableOpenRequests()
	if err != nil {
		return nil, err
	}
	for _, row_data := range tor_data {
		tor.AddRow(row_data["thesisType"], row_data["name"], row_data["thesisTitle"], row_data["requestDate"], row_data["semester"], row_data["status"], "mailto:"+row_data["email"], "/view?tuid="+row_data["tuid"], "#")
	}
	return tor, nil
}

func fillTableMySupervisions(user_id string) (view_protected_homepage.TableMySupervisions, error) {
	tab_data, err := db_model.GetDataTableMySupervisions(user_id)
	if err != nil {
		return nil, err
	}
	tmsv := view_protected_homepage.NewTableMySupervisions()
	for _, row_data := range tab_data {
		tmsv.AddRow(row_data["thesistype"], row_data["name"], row_data["thesistitle"], row_data["deadline"], row_data["supervisor"], row_data["semester"], row_data["thesisstatus"], "mailto:"+row_data["email"], "/view?tuid="+row_data["tuid"], "#")
	}
	return tmsv, nil
}
