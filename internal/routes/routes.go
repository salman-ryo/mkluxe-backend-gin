package routes

import (
	"mkluxe-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter will initialize the Gin engine and mount all route groups
func SetupRouter() *gin.Engine {
	router := gin.Default() //includes default logging and recovery middleware

	// Space for global middleware (CORS, Requst ID, etc.)

	// Register health check route
	router.GET("/health", handler.Health)

	// Route groups to be expanded later
	// public := router.Group("/api/v1")
	// auth := router.Group("/api/v1/auth")
	// admin := router.Group("/api/v1/admin")

	return router

}
