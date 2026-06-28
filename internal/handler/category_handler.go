package handler

import (
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: svc}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	cat, err := h.categoryService.CreateCategory(c.Request.Context(), req.Name, req.Description, req.SortOrder)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Category created successfully", cat)
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to fetch categories")
		return
	}

	response.OK(c, "Categories fetched successfully", categories)
}
