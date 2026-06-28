package handler

import (
	"strconv"

	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"
	"mkluxe-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	prod, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Product created successfully", prod)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "-created_at")

	page, limit, sort = utils.NormalizePagination(page, limit, sort)

	filter := dto.FilterRequest{
		Search:     c.Query("search"),
		CategoryID: c.Query("category_id"),
		Status:     c.Query("status"),
	}

	products, total, err := h.productService.ListProducts(c.Request.Context(), filter, page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to fetch products")
		return
	}

	response.Paginated(c, "Products fetched successfully", products, total, page, limit)
}
