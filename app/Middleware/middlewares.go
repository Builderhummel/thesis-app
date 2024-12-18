package middleware

import (
	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth_controller.TokenValid(c)
		if err != nil {
			c.Redirect(302, "login")
			c.Abort()
			return
		}
		c.Next()
	}
}
