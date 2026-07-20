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

	// 💡 Bind JSON first to get the category slug from the body
	if err := c.ShouldBindJSON(&req); err != nil {
		// 💡 Appended err.Error() to easily debug validation failures
		response.BadRequest(c, "Invalid request payload: "+err.Error(), nil)
		return
	}

	// 💡 Extract from the struct rather than c.Param
	categoryIdentifier := req.CategorySlug
	if categoryIdentifier == "" {
		response.BadRequest(c, "Category identifier is required", nil)
		return
	}

	prod, err := h.productService.CreateProduct(c.Request.Context(), categoryIdentifier, &req)
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

	// Safely parse boolean query parameters into pointers
	var isFeatured, isMostSold *bool
	if featStr := c.Query("is_featured"); featStr != "" {
		val := featStr == "true"
		isFeatured = &val
	}
	if soldStr := c.Query("is_most_sold"); soldStr != "" {
		val := soldStr == "true"
		isMostSold = &val
	}

	filter := dto.FilterRequest{
		Search:     c.Query("search"),
		CategoryID: c.Query("category_id"),
		Status:     c.Query("status"),
		IsFeatured: isFeatured,
		IsMostSold: isMostSold,
	}

	products, total, err := h.productService.ListProducts(c.Request.Context(), filter, page, limit)
	if err != nil {
		// Appending the actual error so you can see why the database or service failed
		response.InternalServerError(c, "Failed to fetch products: "+err.Error())
		return
	}

	response.Paginated(c, "Products fetched successfully", products, total, page, limit)
}

func (h *ProductHandler) Get(c *gin.Context) {
	identifier := c.Param("identifier")

	product, err := h.productService.GetProduct(c.Request.Context(), identifier)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	response.OK(c, "Product fetched successfully", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Product deleted successfully", nil)
}
