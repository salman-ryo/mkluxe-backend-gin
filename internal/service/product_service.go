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

func (s *ProductService) CreateProduct(ctx context.Context, categoryIdentifier string, req *dto.CreateProductRequest) (*domain.Product, error) {
	if err := validation.ValidateProductPayload(req); err != nil {
		return nil, err
	}

	if id, err := primitive.ObjectIDFromHex(categoryIdentifier); err == nil {
		cat, getErr := s.categoryRepo.GetByID(ctx, id)
		if getErr != nil || cat == nil {
			return nil, errors.New("provided primary category ID does not exist")
		}
	} else {
		cat, getErr := s.categoryRepo.GetBySlug(ctx, categoryIdentifier)
		if getErr != nil || cat == nil {
			return nil, errors.New("provided primary category slug does not exist")
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
		return nil, errors.New("a product with this slug already exists")
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
	if len(req.Media) > 0 {
		product.Media = req.Media
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}
	return s.productRepo.Delete(ctx, objID)
}
