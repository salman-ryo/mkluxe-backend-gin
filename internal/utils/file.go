package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func SanitizeFilename(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	base := strings.TrimSuffix(originalName, ext)
	return fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), GenerateSlug(base), ext)
}

func IsAllowedImageExt(filename string) bool {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg", ".png", ".webp", ".svg":
		return true
	}
	return false
}
