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

	//Summary
	summary, err := getSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //TODO: Proper error handling
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
		"Navbar":    renderNavbar(),
		"Summary":   summary,
		"TabOpReq":  tabOpReq,
		"TabMySupV": tabMySupV,
	})
}

func getSummary() (view_protected_homepage.Summary, error) {
	summary_data, err := db_model.GetHomepageRCW()
	if err != nil {
		return view_protected_homepage.Summary{}, err
	}

	summary := view_protected_homepage.NewSummary()

	summary.SetSummary(summary_data["requested"], summary_data["contacted"], summary_data["working"])
	return summary, nil
}
