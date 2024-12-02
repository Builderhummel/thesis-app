package public_controller

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	public := r.Group("/login")
	public.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "public/loginpage/index.html", gin.H{})
	})
	public.POST("/", auth_controller.Login)
}
