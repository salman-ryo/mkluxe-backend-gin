package middleware

import (
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWTs from the access_token cookie and extracts claims
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 💡 Read the token from the secure cookie instead of the Authorization header
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			response.Unauthorized(c, "Authentication cookie is missing or expired")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Attach user context for downstream handlers and RBAC
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
