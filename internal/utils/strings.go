package utils

import "strings"

func CleanString(s string) string {
	return strings.TrimSpace(s)
}

func ToTitleCase(s string) string {
	return strings.Title(strings.ToLower(CleanString(s)))
}
