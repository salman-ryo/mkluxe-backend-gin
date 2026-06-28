package handler

import (
	"fmt"
	"net/http"

	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	mediaService *service.MediaService
}

func NewMediaHandler(svc *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: svc}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file provided", nil)
		return
	}

	safeName, urlPath, err := h.mediaService.ProcessUpload(file.Filename)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// For Phase 6, we are saving it to the local disk.
	// Make sure an 'uploads' directory exists at your project root.
	savePath := fmt.Sprintf("./uploads/%s", safeName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.InternalServerError(c, "Failed to save file to disk")
		return
	}

	response.Success(c, http.StatusOK, "File uploaded successfully", gin.H{
		"url":       urlPath,
		"file_name": safeName,
	})
}
