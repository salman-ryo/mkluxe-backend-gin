package dto

import "mkluxe-backend/internal/domain"

// CreateProductRequest is used for creating new products. All fields are required except for Slug, which can be auto-generated if omitted.
type CreateProductRequest struct {
	Name                string           `json:"name" binding:"required"`
	Slug                string           `json:"slug"` // Auto-generated if omitted
	Description         string           `json:"description" binding:"required"`
	PrimaryCategoryID   string           `json:"primary_category_id" binding:"required"`
	SecondaryCategories []string         `json:"secondary_categories"`
	Status              string           `json:"status" binding:"required,oneof=draft published archived"`
	Variants            []domain.Variant `json:"variants" binding:"required,dive"`
	Media               []domain.Media   `json:"media" binding:"required,dive"`
	FAQs                []domain.FAQ     `json:"faqs"`
	MetaTitle           string           `json:"meta_title"`
	MetaDescription     string           `json:"meta_description"`
}

// UpdateProductRequest is used for updating existing products. All fields are optional, allowing partial updates.
type UpdateProductRequest struct {
	Name                string           `json:"name"`
	Slug                string           `json:"slug"`
	Description         string           `json:"description"`
	PrimaryCategoryID   string           `json:"primary_category_id"`
	SecondaryCategories []string         `json:"secondary_categories"`
	Status              string           `json:"status" binding:"omitempty,oneof=draft published archived"`
	Variants            []domain.Variant `json:"variants"`
	Media               []domain.Media   `json:"media"`
	FAQs                []domain.FAQ     `json:"faqs"`
	MetaTitle           string           `json:"meta_title"`
	MetaDescription     string           `json:"meta_description"`
}

// ProductImportPayload is used for bulk importing products via a JSON payload.
type ProductImportPayload struct {
	Products []CreateProductRequest `json:"products" binding:"required,dive"`
}
