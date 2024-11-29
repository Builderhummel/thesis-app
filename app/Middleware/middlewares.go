package middleware

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth_controller.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized") //TODO: Add login page with error message
			c.Abort()
			return
		}
		c.Next()
	}
}
