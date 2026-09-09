package validation

import (
	"testing"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/dto"
)

func TestValidateProductPayload_SingleVariantAutoDefault(t *testing.T) {
	req := &dto.CreateProductRequest{
		Name:         "Gold Ring",
		Description:  "Luxury ring",
		CategorySlug: "rings",
		Status:       "published",
		Variants: []domain.Variant{
			{
				SKU:       "RNG-001",
				Price:     999.0,
				Stock:     10,
				IsDefault: false,
			},
		},
		Media: []domain.Media{
			{
				URL: "https://example.com/ring.jpg",
			},
		},
	}

	err := ValidateProductPayload(req)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if !req.Variants[0].IsDefault {
		t.Errorf("expected single variant to be auto-marked as default (is_default: true)")
	}

	if !req.Media[0].IsPrimary {
		t.Errorf("expected media to be auto-marked as primary (is_primary: true)")
	}
}

func TestValidateProductPayload_MultipleVariantsAutoDefaultFirst(t *testing.T) {
	req := &dto.CreateProductRequest{
		Name:         "Gold Ring",
		Description:  "Luxury ring",
		CategorySlug: "rings",
		Status:       "published",
		Variants: []domain.Variant{
			{
				SKU:       "RNG-001",
				Price:     999.0,
				Stock:     10,
				IsDefault: false,
			},
			{
				SKU:       "RNG-002",
				Price:     1299.0,
				Stock:     5,
				IsDefault: false,
			},
		},
	}

	err := ValidateProductPayload(req)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if !req.Variants[0].IsDefault {
		t.Errorf("expected first variant to be auto-marked as default")
	}
	if req.Variants[1].IsDefault {
		t.Errorf("expected second variant to remain non-default")
	}
}

func TestValidateProductPayload_MultipleDefaultsError(t *testing.T) {
	req := &dto.CreateProductRequest{
		Name:         "Gold Ring",
		Description:  "Luxury ring",
		CategorySlug: "rings",
		Status:       "published",
		Variants: []domain.Variant{
			{
				SKU:       "RNG-001",
				Price:     999.0,
				Stock:     10,
				IsDefault: true,
			},
			{
				SKU:       "RNG-002",
				Price:     1299.0,
				Stock:     5,
				IsDefault: true,
			},
		},
	}

	err := ValidateProductPayload(req)
	if err == nil {
		t.Fatalf("expected error for multiple default variants, got nil")
	}
}

func TestValidateProductPayload_DuplicateSKU(t *testing.T) {
	req := &dto.CreateProductRequest{
		Name:         "Gold Ring",
		Description:  "Luxury ring",
		CategorySlug: "rings",
		Status:       "published",
		Variants: []domain.Variant{
			{
				SKU:       "RNG-001",
				Price:     999.0,
				Stock:     10,
				IsDefault: true,
			},
			{
				SKU:       "RNG-001",
				Price:     1299.0,
				Stock:     5,
				IsDefault: false,
			},
		},
	}

	err := ValidateProductPayload(req)
	if err == nil {
		t.Fatalf("expected error for duplicate SKU, got nil")
	}
}

func TestValidateProductPayload_InvalidPrice(t *testing.T) {
	req := &dto.CreateProductRequest{
		Name:         "Gold Ring",
		Description:  "Luxury ring",
		CategorySlug: "rings",
		Status:       "published",
		Variants: []domain.Variant{
			{
				SKU:       "RNG-001",
				Price:     0,
				Stock:     10,
				IsDefault: true,
			},
		},
	}

	err := ValidateProductPayload(req)
	if err == nil {
		t.Fatalf("expected error for price <= 0, got nil")
	}
}
