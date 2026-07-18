package routes

import (
	"mkluxe-backend/internal/config"
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

	// CORS Middleware
	router.Use(middleware.CORS(cfg.FrontendURLs))

	// Health check
	router.GET("/health", handler.Health)

	api := router.Group("/api/v1")

	// ==========================================
	// PUBLIC ROUTES (No Auth Required)
	// ==========================================

	// Categories
	api.GET("/categories", handlers.Category.List)
	api.GET("/categories/:identifier", handlers.Category.Get) // identifier can be ID or slug

	// Products
	api.GET("/products", handlers.Product.List)
	api.GET("/products/:identifier", handlers.Product.Get)

	// Inquiries
	api.POST("/inquiries", handlers.Inquiry.Create)

	// Auth
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", handlers.Auth.Login)
		authGroup.POST("/refresh", handlers.Auth.Refresh)
		authGroup.POST("/logout", handlers.Auth.Logout)
		authGroup.GET("/me", middleware.AuthMiddleware(), handlers.Auth.CurrentUser)
	}

	// ==========================================
	// ADMIN ROUTES (Requires JWT Cookie)
	// ==========================================
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	{
		// Categories
		adminGroup.POST("/categories", handlers.Category.Create)
		adminGroup.PUT("/categories/:id", handlers.Category.Update)
		adminGroup.DELETE("/categories/:id", handlers.Category.Delete)

		// Products
		adminGroup.POST("/products", handlers.Product.Create) // 💡 Updated to /products
		adminGroup.PUT("/products/:id", handlers.Product.Update)
		adminGroup.DELETE("/products/:id", handlers.Product.Delete)

		// Inquiries
		// adminGroup.GET("/inquiries", handlers.Inquiry.List)
		// adminGroup.GET("/inquiries/:id", handlers.Inquiry.Get)
		adminGroup.PATCH("/inquiries/:id/status", handlers.Inquiry.UpdateStatus)
		// adminGroup.DELETE("/inquiries/:id", handlers.Inquiry.Delete)
	}

	return router
}
