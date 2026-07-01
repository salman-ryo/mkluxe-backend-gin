package service

import (
	"context"
	"errors"

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
}

func NewProductService(pRepo *repository.ProductRepository, cRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{productRepo: pRepo, categoryRepo: cRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*domain.Product, error) {
	if err := validation.ValidateProductPayload(req); err != nil {
		return nil, err
	}

	// Verify Primary Category exists
	primaryCatID, _ := primitive.ObjectIDFromHex(req.PrimaryCategoryID)
	_, err := s.categoryRepo.GetByID(ctx, primaryCatID)
	if err != nil {
		return nil, errors.New("provided primary category does not exist")
	}

	// Normalize Slug
	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Name)
	} else {
		slug = utils.GenerateSlug(slug)
	}

	existing, _ := s.productRepo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("a product with this slug already exists")
	}

	// Convert Secondary Category strings to ObjectIDs safely
	var secCatIDs []primitive.ObjectID
	for _, idStr := range req.SecondaryCategories {
		if id, err := primitive.ObjectIDFromHex(idStr); err == nil {
			secCatIDs = append(secCatIDs, id)
		}
	}

	product := &domain.Product{
		Name:                utils.CleanString(req.Name),
		Slug:                slug,
		Description:         utils.CleanString(req.Description),
		PrimaryCategoryID:   primaryCatID,
		SecondaryCategories: secCatIDs,
		Status:              req.Status,

		// Mapping the new flag values
		IsFeatured: req.IsFeatured,
		IsMostSold: req.IsMostSold,

		Variants:        req.Variants,
		Media:           req.Media,
		FAQs:            req.FAQs,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, filter dto.FilterRequest, page, limit int) ([]domain.Product, int64, error) {
	return s.productRepo.List(ctx, filter, page, limit)
}
