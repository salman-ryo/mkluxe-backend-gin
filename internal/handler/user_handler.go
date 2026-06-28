package handler

import (
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{userService: svc}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Admin user created successfully", user)
}
