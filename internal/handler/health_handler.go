package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Checks server health
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Server is up and running",
		"status":  http.StatusOK,
	})
}
