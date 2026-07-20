package middleware

import (
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWTs from the access_token cookie and extracts claims.
// It features automatic access token renewal if the access token is expired
// but a valid refresh_token is present.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *utils.JWTClaims
		var err error
		needsRefresh := false

		// 💡 Read the access token from the secure cookie
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			needsRefresh = true
		} else {
			claims, err = utils.ValidateToken(tokenString) //[cite: 1]
			if err != nil {
				// Access token is invalid or expired, trigger refresh
				needsRefresh = true
			}
		}

		if needsRefresh {
			// 1. Retrieve the refresh token
			refreshTokenString, err := c.Cookie("refresh_token")
			if err != nil || refreshTokenString == "" {
				response.Unauthorized(c, "Authentication cookie is missing or expired") //[cite: 2]
				c.Abort()
				return
			}

			// 2. Validate the refresh token
			refreshClaims, err := utils.ValidateToken(refreshTokenString) //[cite: 1]
			if err != nil {
				response.Unauthorized(c, "Session expired, please log in again") //[cite: 2]
				c.Abort()
				return
			}

			// 3. Generate a new access token using claims from the valid refresh token
			// We discard the newly generated refresh token because we only want to renew the access token.
			newAccessToken, _, err := utils.GenerateTokens(refreshClaims.UserID, refreshClaims.Role) //[cite: 1]
			if err != nil {
				response.InternalServerError(c, "Failed to renew access token") //[cite: 2]
				c.Abort()
				return
			}

			// 4. Set the new access token as an HTTP-only cookie
			// MaxAge: 900 seconds (15 minutes) matches your access token lifetime[cite: 1]
			// Path: "/"
			// Domain: "" (defaults to current domain)
			// Secure: false (Set this to true in production with HTTPS)
			// HttpOnly: true
			c.SetCookie("access_token", newAccessToken, 900, "/", "", false, true)

			// 5. Assign the validated claims from the refresh token to the current request
			claims = refreshClaims
		}

		// Attach user context for downstream handlers and RBAC
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
