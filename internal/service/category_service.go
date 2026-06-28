package service

import (
	"context"
	"errors"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, name, description string, sortOrder int) (*domain.Category, error) {
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
		SortOrder:   sortOrder,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.ListAll(ctx)
}
