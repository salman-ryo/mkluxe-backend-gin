package middleware

import (
	"mkluxe-backend/internal/constants"
	"mkluxe-backend/internal/response"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware restricts route access to specific user roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			response.Unauthorized(c, "User role not found in context")
			c.Abort()
			return
		}

		roleStr := userRole.(string)

		// Super Admin always bypasses role checks
		if roleStr == constants.RoleSuperAdmin {
			c.Next()
			return
		}

		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "You do not have permission to access this resource")
		c.Abort()
	}
}
