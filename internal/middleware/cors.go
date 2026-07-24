package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS takes the allowed frontend URLs from our config layer
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := strings.TrimSuffix(c.Request.Header.Get("Origin"), "/")

		// Check if the incoming origin is in our allowed list
		isAllowed := false
		for _, o := range allowedOrigins {
			if origin == strings.TrimSuffix(o, "/") {
				isAllowed = true
				break
			}
		}

		// If it's a match, echo that specific origin back (required for cookies)
		if isAllowed {
			// Set the explicit matched origin to allow credentials
			c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
