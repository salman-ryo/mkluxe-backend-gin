package handler

import (
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	r2Service *service.R2Service
}

func NewUploadHandler(r2 *service.R2Service) *UploadHandler {
	return &UploadHandler{r2Service: r2}
}

type PresignedUploadRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	ObjectKey   string `json:"object_key"`
}

func (h *UploadHandler) GetPresignedURL(c *gin.Context) {
	var req PresignedUploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload: file_name and content_type are required", nil)
		return
	}

	expires := 15 * time.Minute

	uploadURL, publicURL, key, err := h.r2Service.GetPresignedUploadURL(
		c.Request.Context(),
		req.FileName,
		req.ContentType,
		req.ObjectKey,
		expires,
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "Presigned upload URL generated successfully", gin.H{
		"upload_url": uploadURL,
		"public_url": publicURL,
		"key":        key,
	})
}
