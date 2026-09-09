package service

import (
	"testing"
)

func TestExtractObjectKey(t *testing.T) {
	svc := &R2Service{
		publicBaseURL: "https://pub-abcdef123456.r2.dev",
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard public base URL matching",
			input:    "https://pub-abcdef123456.r2.dev/uploads/1725890000000000000_ring.webp",
			expected: "uploads/1725890000000000000_ring.webp",
		},
		{
			name:     "URL with query params",
			input:    "https://pub-abcdef123456.r2.dev/uploads/1725890000000000000_ring.webp?width=400&format=webp",
			expected: "uploads/1725890000000000000_ring.webp",
		},
		{
			name:     "generic r2.dev domain matching",
			input:    "https://custom-bucket.r2.dev/uploads/necklace_gold.png",
			expected: "uploads/necklace_gold.png",
		},
		{
			name:     "r2 cloudflarestorage endpoint",
			input:    "https://account123.r2.cloudflarestorage.com/my-bucket/uploads/earrings.jpg",
			expected: "uploads/earrings.jpg",
		},
		{
			name:     "relative path starting with slash",
			input:    "/uploads/1234_bracelet.png",
			expected: "uploads/1234_bracelet.png",
		},
		{
			name:     "clean key already",
			input:    "uploads/1234_bracelet.png",
			expected: "uploads/1234_bracelet.png",
		},
		{
			name:     "external URL - unsplash",
			input:    "https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?w=500",
			expected: "",
		},
		{
			name:     "external URL - cloudinary",
			input:    "https://res.cloudinary.com/demo/image/upload/sample.jpg",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "path traversal attempt",
			input:    "https://pub-abcdef123456.r2.dev/uploads/../../etc/passwd",
			expected: "",
		},
		{
			name:     "just uploads/ directory prefix without filename",
			input:    "https://pub-abcdef123456.r2.dev/uploads/",
			expected: "",
		},
		{
			name:     "hierarchical products path",
			input:    "https://pub-abcdef123456.r2.dev/products/rings/diamond-solitaire/1725893000000_1_front.webp",
			expected: "products/rings/diamond-solitaire/1725893000000_1_front.webp",
		},
		{
			name:     "hierarchical products path with query params",
			input:    "https://pub-abcdef123456.r2.dev/products/necklaces/gold-chain/1725893000000_2_detail.jpg?width=600",
			expected: "products/necklaces/gold-chain/1725893000000_2_detail.jpg",
		},
		{
			name:     "category path",
			input:    "https://pub-abcdef123456.r2.dev/categories/rings.webp",
			expected: "categories/rings.webp",
		},
		{
			name:     "generic r2 products endpoint",
			input:    "https://my-bucket.r2.dev/products/bracelets/silver-cuff/123_main.png",
			expected: "products/bracelets/silver-cuff/123_main.png",
		},
		{
			name:     "relative products path",
			input:    "/products/rings/vintage-band/456.webp",
			expected: "products/rings/vintage-band/456.webp",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := svc.ExtractObjectKey(tc.input)
			if actual != tc.expected {
				t.Errorf("ExtractObjectKey(%q) = %q, want %q", tc.input, actual, tc.expected)
			}
		})
	}
}
