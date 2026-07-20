package dto

import "mkluxe-backend/internal/domain"

// CreateProductRequest is used for creating new products.
type CreateProductRequest struct {
	CategorySlug    string           `json:"category_slug" binding:"required"` // 💡 Added this to read from body
	Name            string           `json:"name" binding:"required"`
	Slug            string           `json:"slug"` // Auto-generated if omitted
	Description     string           `json:"description" binding:"required"`
	Status          string           `json:"status" binding:"required,oneof=draft published archived"`
	IsFeatured      bool             `json:"is_featured"`
	IsMostSold      bool             `json:"is_most_sold"`
	Variants        []domain.Variant `json:"variants" binding:"required,dive"`
	Media           []domain.Media   `json:"media" binding:"required,dive"`
	FAQs            []domain.FAQ     `json:"faqs"`
	MetaTitle       string           `json:"meta_title"`
	MetaDescription string           `json:"meta_description"`
}

// UpdateProductRequest remains unchanged
type UpdateProductRequest struct {
	Name            string           `json:"name"`
	Slug            string           `json:"slug"`
	Description     string           `json:"description"`
	CategorySlug    string           `json:"category_slug"`
	Status          string           `json:"status" binding:"omitempty,oneof=draft published archived"`
	IsFeatured      *bool            `json:"is_featured"`
	IsMostSold      *bool            `json:"is_most_sold"`
	Variants        []domain.Variant `json:"variants"`
	Media           []domain.Media   `json:"media"`
	FAQs            []domain.FAQ     `json:"faqs"`
	MetaTitle       string           `json:"meta_title"`
	MetaDescription string           `json:"meta_description"`
}

type ProductImportPayload struct {
	Products []CreateProductRequest `json:"products" binding:"required,dive"`
}
