package handler

import (
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: svc}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve dashboard statistics")
		return
	}

	response.OK(c, "Dashboard statistics fetched successfully", stats)
}
