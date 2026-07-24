package middleware

import (
	"mkluxe-backend/internal/config"
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWTs from the access_token cookie and extracts claims.
// It features automatic access token renewal if the access token is expired
// but a valid refresh_token is present.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *utils.JWTClaims
		var err error
		needsRefresh := false

		accessCfg := cfg.GetAccessCookieConfig()
		refreshCfg := cfg.GetRefreshCookieConfig()

		// 💡 Read the access token from the secure cookie
		tokenString, err := c.Cookie(accessCfg.Name)
		if err != nil || tokenString == "" {
			needsRefresh = true
		} else {
			claims, err = utils.ValidateToken(tokenString, cfg.JWTAccessSecret)
			if err != nil {
				// Access token is invalid or expired, trigger refresh
				needsRefresh = true
			}
		}

		if needsRefresh {
			// 1. Retrieve the refresh token
			refreshTokenString, err := c.Cookie(refreshCfg.Name)
			if err != nil || refreshTokenString == "" {
				response.Unauthorized(c, "Authentication cookie is missing or expired")
				c.Abort()
				return
			}

			// 2. Validate the refresh token
			refreshClaims, err := utils.ValidateToken(refreshTokenString, cfg.JWTRefreshSecret)
			if err != nil {
				response.Unauthorized(c, "Session expired, please log in again")
				c.Abort()
				return
			}

			// 3. Generate a new access token using claims from the valid refresh token
			// We discard the newly generated refresh token because we only want to renew the access token.
			newAccessToken, _, err := utils.GenerateTokens(refreshClaims.UserID, refreshClaims.Role, cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
			if err != nil {
				response.InternalServerError(c, "Failed to renew access token")
				c.Abort()
				return
			}

			// 4. Set the new access token as an HTTP-only cookie using parameters from config
			c.SetCookie(accessCfg.Name, newAccessToken, accessCfg.MaxAge, accessCfg.Path, accessCfg.Domain, accessCfg.Secure, accessCfg.HttpOnly)

			// 5. Assign the validated claims from the refresh token to the current request
			claims = refreshClaims
		}

		// Attach user context for downstream handlers and RBAC
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
