package handler

import (
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type InquiryHandler struct {
	inquiryService *service.InquiryService
}

func NewInquiryHandler(svc *service.InquiryService) *InquiryHandler {
	return &InquiryHandler{inquiryService: svc}
}

func (h *InquiryHandler) Create(c *gin.Context) {
	var req dto.CreateInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	inq, err := h.inquiryService.CreateInquiry(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Inquiry submitted successfully", inq)
}

func (h *InquiryHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateInquiryStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	if err := h.inquiryService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Inquiry status updated successfully", nil)
}
