package middleware

import (
	"net/http"

	"github.com/Builderhummel/thesis-app/app/Constants/roles"
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

func RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth_controller.GetUserRoles(c)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRole(requiredRoles ...roles.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIface, _ := c.Get("role")
		userRole, _ := roleIface.(roles.Role)
		for _, r := range requiredRoles {
			if userRole == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}
