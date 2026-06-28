package handler

import (
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: svc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	authRes, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "Login successful", authRes)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	authRes, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "Token refreshed successfully", authRes)
}

func (h *AuthHandler) CurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized request")
		return
	}

	userRes, err := h.authService.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "Current user fetched successfully", userRes)
}
