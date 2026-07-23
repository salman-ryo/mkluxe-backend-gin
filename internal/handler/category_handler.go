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
		IsFeatured  bool   `json:"is_featured"` // 💡 Added here
		ImageURL    string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	cat, err := h.categoryService.CreateCategory(c.Request.Context(), req.Name, req.Description, req.SortOrder, req.IsFeatured, req.ImageURL)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Category created successfully", cat)
}

func (h *CategoryHandler) List(c *gin.Context) {
	// 💡 Read the query parameter out of the URL (e.g., ?is_featured=true)
	var isFeatured *bool
	if featStr := c.Query("is_featured"); featStr != "" {
		val := featStr == "true"
		isFeatured = &val
	}

	categories, err := h.categoryService.ListCategories(c.Request.Context(), isFeatured)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch categories")
		return
	}

	response.OK(c, "Categories fetched successfully", categories)
}

func (h *CategoryHandler) Get(c *gin.Context) {
	identifier := c.Param("identifier")

	cat, err := h.categoryService.GetCategory(c.Request.Context(), identifier)
	if err != nil {
		response.NotFound(c, "Category not found")
		return
	}

	response.OK(c, "Category fetched successfully", cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		SortOrder   *int    `json:"sort_order"`
		IsActive    *bool   `json:"is_active"`
		IsFeatured  *bool   `json:"is_featured"` // 💡 Added here
		ImageURL    *string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	cat, err := h.categoryService.UpdateCategory(c.Request.Context(), id, req.Name, req.Description, req.SortOrder, req.IsActive, req.IsFeatured, req.ImageURL)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Category updated successfully", cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Category deleted successfully", nil)
}
