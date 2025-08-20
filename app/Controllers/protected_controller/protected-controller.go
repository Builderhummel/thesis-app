package protected_controller

import (
	middleware "github.com/Builderhummel/thesis-app/app/Middleware"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("/", RenderHomepage)

	protected.GET("/all", RenderAllSupervisions)
	protected.StaticFile("/js/custom/all/datatables.js", "./app/Views/templates/protected/all_supervisions/datatables.js")

	protected.GET("/my_supervisions", RenderMySupervisions)
	protected.StaticFile("/js/custom/my_supervisions/datatables.js", "./app/Views/templates/protected/my_supervisions/datatables.js")

	protected.GET("/open_requests", RenderOpenRequests)
	protected.StaticFile("/js/custom/open_requests/datatables.js", "./app/Views/templates/protected/open_requests/datatables.js")

	protected.GET("/add", RenderAddSupervisionRequestForm)
	protected.POST("/add", HandleNewSupervisionRequest)

	protected.GET("/view", RenderViewSupervisionRequestForm)

	protected.GET("/edit", RenderEditSupervisionRequestForm)
	protected.POST("/edit", HandleEditSupervisionRequest)
	protected.StaticFile("/js/custom/edit/select2.js", "./app/Views/templates/protected/edit_supervison_request/select2.js")

	protected.GET("/users", RenderListAllUsers)

	protected.GET("/new_user", RenderNewUser)
	protected.POST("/new_user", HandlePostNewUser)

	protected.GET("/edit_user", RenderEditUser)
	protected.POST("/edit_user", HandlePostEditUser)

	protected.GET("/export", HandleExport)
}
