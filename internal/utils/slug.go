package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug converts a string into a URL-friendly slug by lowercasing it, replacing non-alphanumeric characters with hyphens, and trimming extra hyphens.
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func GenerateSlug(input string) string {
	// Normalize input: trim whitespace and convert to lowercase
	cleaned := strings.ToLower(strings.TrimSpace(input))

	// Replace any sequence of non-alphanumeric characters with a hyphen
	slug := nonAlphanumericRegex.ReplaceAllString(cleaned, "-")

	// Remove leading/trailing hyphens
	return strings.Trim(slug, "-")
}
