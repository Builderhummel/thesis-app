package auth_controller

import (
	"github.com/Builderhummel/thesis-app/app/Constants/roles"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/gin-gonic/gin"
)

func GetUserRoles(c *gin.Context) error {
	userID := c.GetString("user_id")
	role, err := db_model.GetUserRoleByLoginHandle(userID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "failed to fetch role"})
		return err
	}
	c.Set("role", role)
	return nil
}

func GetUserRoleFromContext(c *gin.Context) roles.Role {
	if val, exists := c.Get("role"); exists {
		if cRole, ok := val.(roles.Role); ok {
			return cRole
		}
		if roleStr, ok := val.(string); ok {
			return roles.Role(roleStr)
		}
	}
	return roles.RoleDefault
}
