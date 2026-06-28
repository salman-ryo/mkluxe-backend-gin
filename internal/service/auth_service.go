package service

import (
	"context"
	"errors"
	"time"

	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// this method belongs to AuthService, s stands for 'self'
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, err
	}

	// Update last login
	user.LastLogin = time.Now()
	_ = s.userRepo.Update(ctx, user)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID.Hex(),
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}
