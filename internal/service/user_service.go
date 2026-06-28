// Manages admin account lifecycles and safely hashes passwords.

package service

import (
	"context"
	"errors"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"
	"mkluxe-backend/internal/validation"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*domain.User, error) {
	if err := validation.ValidateAdminUser(req); err != nil {
		return nil, err
	}

	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email is already in use")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
