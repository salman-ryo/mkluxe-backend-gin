package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"
	"mkluxe-backend/internal/validation"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	r2Service    *R2Service
}

func NewProductService(pRepo *repository.ProductRepository, cRepo *repository.CategoryRepository, r2Svc *R2Service) *ProductService {
	return &ProductService{
		productRepo:  pRepo,
		categoryRepo: cRepo,
		r2Service:    r2Svc,
	}
}

// ValidateProduct checks business rules, category existence, and slug uniqueness without creating the product
func (s *ProductService) ValidateProduct(ctx context.Context, categoryIdentifier string, req *dto.CreateProductRequest, excludeProductID string) (string, error) {
	if err := validation.ValidateProductPayload(req); err != nil {
		return "", err
	}

	if categoryIdentifier == "" {
		categoryIdentifier = req.CategorySlug
	}
	if categoryIdentifier == "" {
		return "", errors.New("category identifier is required")
	}

	if id, err := primitive.ObjectIDFromHex(categoryIdentifier); err == nil {
		cat, getErr := s.categoryRepo.GetByID(ctx, id)
		if getErr != nil || cat == nil {
			return "", errors.New("provided primary category ID does not exist")
		}
	} else {
		cat, getErr := s.categoryRepo.GetBySlug(ctx, categoryIdentifier)
		if getErr != nil || cat == nil {
			return "", errors.New("provided primary category slug does not exist")
		}
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	} else {
		slug = utils.GenerateSlug(slug)
	}

	existing, _ := s.productRepo.GetBySlug(ctx, slug)
	if existing != nil {
		if excludeProductID == "" || existing.ID.Hex() != excludeProductID {
			return "", errors.New("a product with this slug already exists")
		}
	}

	return slug, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, categoryIdentifier string, req *dto.CreateProductRequest) (*domain.Product, error) {
	slug, err := s.ValidateProduct(ctx, categoryIdentifier, req, "")
	if err != nil {
		return nil, err
	}

	metaTitle := req.MetaTitle
	if metaTitle == "" {
		metaTitle = fmt.Sprintf("%s | MK Luxe Divine", utils.CleanString(req.Name))
	}

	metaDescription := req.MetaDescription
	if metaDescription == "" {
		metaDescription = utils.CleanString(req.Description)
		if len(metaDescription) > 155 {
			metaDescription = metaDescription[:152] + "..."
		}
	}

	// Normalize media alt_text if alt was provided
	for i := range req.Media {
		if req.Media[i].AltText == "" && req.Media[i].Alt != "" {
			req.Media[i].AltText = req.Media[i].Alt
		}
	}

	product := &domain.Product{
		Name:            utils.CleanString(req.Name),
		Slug:            slug,
		Description:     utils.CleanString(req.Description),
		CategorySlug:    categoryIdentifier, // Store the slug for easier querying
		Status:          req.Status,
		IsFeatured:      req.IsFeatured,
		IsMostSold:      req.IsMostSold,
		Variants:        req.Variants,
		Media:           req.Media,
		FAQs:            req.FAQs,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, filter dto.FilterRequest, page, limit int) ([]domain.Product, int64, error) {
	return s.productRepo.List(ctx, filter, page, limit)
}

func (s *ProductService) GetProduct(ctx context.Context, identifier string) (*domain.Product, error) {
	if id, err := primitive.ObjectIDFromHex(identifier); err == nil {
		return s.productRepo.GetByID(ctx, id)
	}
	return s.productRepo.GetBySlug(ctx, identifier)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, req *dto.UpdateProductRequest) (*domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	product, err := s.productRepo.GetByID(ctx, objID)
	if err != nil || product == nil {
		return nil, errors.New("product not found")
	}

	// Apply updates if fields are provided
	if req.Name != "" {
		product.Name = utils.CleanString(req.Name)
	}
	if req.Slug != "" {
		product.Slug = utils.GenerateSlug(req.Slug)
	}
	if req.Description != "" {
		product.Description = utils.CleanString(req.Description)
	}
	if req.CategorySlug != "" {
		product.CategorySlug = req.CategorySlug
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if req.IsMostSold != nil {
		product.IsMostSold = *req.IsMostSold
	}
	if len(req.Variants) > 0 {
		product.Variants = req.Variants
	}

	var removedMediaURLs []string
	if req.Media != nil {
		newURLSet := make(map[string]struct{}, len(req.Media))
		for i := range req.Media {
			if req.Media[i].AltText == "" && req.Media[i].Alt != "" {
				req.Media[i].AltText = req.Media[i].Alt
			}
			if req.Media[i].URL != "" {
				newURLSet[req.Media[i].URL] = struct{}{}
			}
		}

		for _, m := range product.Media {
			if m.URL != "" {
				if _, exists := newURLSet[m.URL]; !exists {
					removedMediaURLs = append(removedMediaURLs, m.URL)
				}
			}
		}

		product.Media = req.Media
	}

	if req.FAQs != nil {
		product.FAQs = req.FAQs
	}
	if req.MetaTitle != "" {
		product.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		product.MetaDescription = req.MetaDescription
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	// Asynchronously delete removed images from R2
	if s.r2Service != nil && len(removedMediaURLs) > 0 {
		go func(urls []string) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := s.r2Service.DeleteMediaByURLs(bgCtx, urls); err != nil {
				log.Printf("[ProductService] Warning: failed to delete removed R2 media on product update (%s): %v", id, err)
			}
		}(removedMediaURLs)
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}

	// 1. Fetch product to get media URLs before deleting
	product, err := s.productRepo.GetByID(ctx, objID)
	if err != nil || product == nil {
		return errors.New("product not found")
	}

	// 2. Collect media URLs to delete
	var mediaURLs []string
	for _, m := range product.Media {
		if m.URL != "" {
			mediaURLs = append(mediaURLs, m.URL)
		}
	}

	// 3. Delete from database
	if err := s.productRepo.Delete(ctx, objID); err != nil {
		return err
	}

	// 4. Asynchronously delete R2 media files
	if s.r2Service != nil && len(mediaURLs) > 0 {
		go func(urls []string) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := s.r2Service.DeleteMediaByURLs(bgCtx, urls); err != nil {
				log.Printf("[ProductService] Warning: failed to delete R2 media on product deletion (%s): %v", id, err)
			}
		}(mediaURLs)
	}

	return nil
}
