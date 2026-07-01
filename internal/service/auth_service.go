package service

import (
	"context"
	"errors"
	"log"
	"time"

	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err != nil {
		log.Printf("Debug - Login failed: GetByEmail error: %v", err)
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		log.Printf("Debug - Login failed: User not found for email: %s", req.Email)
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		log.Printf("Debug - Login failed: User account is disabled")
		return nil, errors.New("account is disabled")
	}

	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		log.Printf("Debug - Login failed: Password verification failed for %s = %s", req.Email, req.Password)
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

// RefreshToken issues a new token pair if the provided refresh token is valid
func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetCurrentUser fetches the authenticated user's profile data
func (s *AuthService) GetCurrentUser(ctx context.Context, userIDStr string) (*dto.UserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}
