package service

import (
	"errors"
	"fmt"

	"mkluxe-backend/internal/utils"
)

type MediaService struct{}

func NewMediaService() *MediaService {
	return &MediaService{}
}

// ProcessUpload normalizes filenames and checks extensions before saving
func (s *MediaService) ProcessUpload(originalFilename string) (string, string, error) {
	if !utils.IsAllowedImageExt(originalFilename) {
		return "", "", errors.New("invalid file type; only images are allowed")
	}

	safeName := utils.SanitizeFilename(originalFilename)
	urlPath := fmt.Sprintf("/uploads/%s", safeName) // Path for frontend consumption

	return safeName, urlPath, nil
}
