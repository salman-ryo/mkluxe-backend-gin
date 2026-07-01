package routes

import (
	"mkluxe-backend/internal/config" // 👈 Added config import
	"mkluxe-backend/internal/handler"
	"mkluxe-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	Auth     *handler.AuthHandler
	Category *handler.CategoryHandler
	Product  *handler.ProductHandler
	Inquiry  *handler.InquiryHandler
}

// SetupRouter initializes Gin and mounts all routes
func SetupRouter(handlers AppHandlers, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// 💡 Pass the dynamically loaded Frontend URLs to the middleware
	router.Use(middleware.CORS(cfg.FrontendURLs))

	// Health check
	router.GET("/health", handler.Health)

	api := router.Group("/api/v1")

	// Mount Public Routes (No Auth Required)
	api.GET("/categories", handlers.Category.List)
	api.GET("/products", handlers.Product.List)
	api.POST("/inquiries", handlers.Inquiry.Create)

	// Mount Auth Routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", handlers.Auth.Login)
		authGroup.POST("/refresh", handlers.Auth.Refresh)
		authGroup.POST("/logout", handlers.Auth.Logout) // 💡 Added your new Logout route!
		authGroup.GET("/me", middleware.AuthMiddleware(), handlers.Auth.CurrentUser)
	}

	// Mount Protected Admin Routes
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware()) // Everything here requires a JWT cookie
	{
		// Categories (Requires Auth)
		adminGroup.POST("/categories", handlers.Category.Create)

		// Products (Requires Auth)
		adminGroup.POST("/products", handlers.Product.Create)

		// Inquiries Admin (Requires Auth)
		adminGroup.PATCH("/inquiries/:id/status", handlers.Inquiry.UpdateStatus)
	}

	return router
}
