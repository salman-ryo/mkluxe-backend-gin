package utils

import (
	"os"
)

// GetEnv handles default values easily (Capital G so it can be exported)
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
