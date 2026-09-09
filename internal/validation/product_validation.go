package validation

import (
	"errors"
	"fmt"
	"mkluxe-backend/internal/dto"
)

func ValidateProductPayload(req *dto.CreateProductRequest) error {
	// 💡 REMOVED the Primary Category ID format check.
	// Since you are passing the category_slug now, req.PrimaryCategoryID is
	// empty at this stage. The Service layer will populate it later.

	// 2. Variant business logic
	if len(req.Variants) == 0 {
		return errors.New("product must include at least one variant")
	}

	// Auto-default if no default is explicitly marked
	defaultCount := 0
	for _, v := range req.Variants {
		if v.IsDefault {
			defaultCount++
		}
	}
	if defaultCount == 0 {
		// Automatically mark the first variant as default
		req.Variants[0].IsDefault = true
		defaultCount = 1
	}
	if defaultCount > 1 {
		return errors.New("multiple variants marked as default; only one allowed")
	}

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
			return fmt.Errorf("variant %s has an invalid price (must be greater than 0)", v.SKU)
		}
	}

	// 3. Media primary image check
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
			// Automatically mark the first media item as primary
			req.Media[0].IsPrimary = true
		}
	}

	return nil
}
