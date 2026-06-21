package validation

import (
	"errors"
	"fmt"
	"mkluxe-backend/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateProductPayload(req *dto.CreateProductRequest) error {
	// 1. Check Primary Category ID format
	if _, err := primitive.ObjectIDFromHex(req.PrimaryCategoryID); err != nil {
		return errors.New("invalid primary_category_id format")
	}

	// 2. Check Secondary Category IDs
	for _, catID := range req.SecondaryCategories {
		if _, err := primitive.ObjectIDFromHex(catID); err != nil {
			return fmt.Errorf("invalid secondary category ID format: %s", catID)
		}
	}

	// 3. Variant business logic
	if len(req.Variants) == 0 {
		return errors.New("product must include at least one variant")
	}

	defaultCount := 0
	skuSet := make(map[string]bool)

	for i, v := range req.Variants {
		if v.SKU == "" {
			return fmt.Errorf("variant at index %d is missing a SKU", i)
		}
		if skuSet[v.SKU] {
			return fmt.Errorf("duplicate SKU found in variant list: %s", v.SKU)
		}
		skuSet[v.SKU] = true

		if v.Price <= 0 {
			return fmt.Errorf("variant %s has an invalid price", v.SKU)
		}
		if v.IsDefault {
			defaultCount++
		}
	}

	if defaultCount == 0 {
		return errors.New("exactly one variant must be marked as default (is_default: true)")
	}
	if defaultCount > 1 {
		return errors.New("multiple variants marked as default; only one allowed")
	}

	// 4. Media primary image check
	if len(req.Media) > 0 {
		primaryCount := 0
		for _, m := range req.Media {
			if m.URL == "" {
				return errors.New("media entry missing URL")
			}
			if m.IsPrimary {
				primaryCount++
			}
		}
		if primaryCount == 0 {
			return errors.New("at least one media image must be marked as primary (is_primary: true)")
		}
	}

	return nil
}
