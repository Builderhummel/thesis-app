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
}
