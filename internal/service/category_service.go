package service

import (
	"context"
	"errors"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// 💡 Added isFeatured parameter
func (s *CategoryService) CreateCategory(ctx context.Context, name, description string, sortOrder int, isFeatured bool) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	slug := utils.GenerateSlug(name)
	existing, _ := s.repo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("a category with this name already exists")
	}

	cat := &domain.Category{
		Name:        utils.ToTitleCase(name),
		Slug:        slug,
		Description: utils.CleanString(description),
		IsActive:    true,
		IsFeatured:  isFeatured, // 💡 Assigned here
		SortOrder:   sortOrder,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// 💡 Accepts the filter parameter and passes it to the repo
func (s *CategoryService) ListCategories(ctx context.Context, isFeatured *bool) ([]domain.Category, error) {
	return s.repo.ListAll(ctx, isFeatured)
}

func (s *CategoryService) GetCategory(ctx context.Context, identifier string) (*domain.Category, error) {
	if id, err := primitive.ObjectIDFromHex(identifier); err == nil {
		return s.repo.GetByID(ctx, id)
	}
	return s.repo.GetBySlug(ctx, identifier)
}

// 💡 Added isFeatured parameter as a pointer
func (s *CategoryService) UpdateCategory(ctx context.Context, id string, name, description *string, sortOrder *int, isActive, isFeatured *bool) (*domain.Category, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	cat, err := s.repo.GetByID(ctx, objID)
	if err != nil || cat == nil {
		return nil, errors.New("category not found")
	}

	if name != nil && *name != "" {
		cat.Name = utils.ToTitleCase(*name)
		cat.Slug = utils.GenerateSlug(*name)
	}
	if description != nil {
		cat.Description = utils.CleanString(*description)
	}
	if sortOrder != nil {
		cat.SortOrder = *sortOrder
	}
	if isActive != nil {
		cat.IsActive = *isActive
	}
	if isFeatured != nil {
		cat.IsFeatured = *isFeatured // 💡 Update handling
	}

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID format")
	}
	return s.repo.Delete(ctx, objID)
}
