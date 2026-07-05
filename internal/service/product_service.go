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

// 💡 Added categoryIdentifier as an argument
func (s *ProductService) CreateProduct(ctx context.Context, categoryIdentifier string, req *dto.CreateProductRequest) (*domain.Product, error) {
	if err := validation.ValidateProductPayload(req); err != nil {
		return nil, err
	}

	var primaryCatID primitive.ObjectID

	// 1. Check if the identifier is a standard MongoDB ID
	if id, err := primitive.ObjectIDFromHex(categoryIdentifier); err == nil {
		cat, getErr := s.categoryRepo.GetByID(ctx, id)
		if getErr != nil || cat == nil {
			return nil, errors.New("provided primary category ID does not exist")
		}
		primaryCatID = cat.ID
	} else {
		// 2. If it's not a Mongo ID, treat it as a Slug
		cat, getErr := s.categoryRepo.GetBySlug(ctx, categoryIdentifier) // ⚠️ Make sure GetBySlug exists in your CategoryRepo!
		if getErr != nil || cat == nil {
			return nil, errors.New("provided primary category slug does not exist")
		}
		primaryCatID = cat.ID
	}

	// Normalize Slug for the Product
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
		PrimaryCategoryID:   primaryCatID, // 💡 Assigned dynamically from the URL check
		SecondaryCategories: secCatIDs,
		Status:              req.Status,
		IsFeatured:          req.IsFeatured,
		IsMostSold:          req.IsMostSold,
		Variants:            req.Variants,
		Media:               req.Media,
		FAQs:                req.FAQs,
		MetaTitle:           req.MetaTitle,
		MetaDescription:     req.MetaDescription,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, filter dto.FilterRequest, page, limit int) ([]domain.Product, int64, error) {
	return s.productRepo.List(ctx, filter, page, limit)
}
