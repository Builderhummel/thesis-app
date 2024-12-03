package protected_controller

import (
	middleware "github.com/Builderhummel/thesis-app/app/Middleware"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("/", RenderHomepage)

	protected.GET("/add", RenderAddSupervisionRequestForm)
	protected.POST("/add", HandleNewSupervisionRequest)
}
