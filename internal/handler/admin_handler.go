package handler

import (
	"mkluxe-backend/internal/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	// Add multiple services here later (ProductService, InquiryService) for stats
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	// Dummy data until Phase 9 is implemented
	stats := gin.H{
		"total_products":    0,
		"pending_inquiries": 0,
	}
	response.OK(c, "Dashboard stats fetched successfully", stats)
}
