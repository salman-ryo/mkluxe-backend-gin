package seed

import (
	"context"
	"log"

	"mkluxe-backend/internal/constants"
	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

// SeedSuperAdmin creates the default admin if no users exist
func SeedSuperAdmin(db *mongo.Database) error {
	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()

	// Check if admin already exists
	existing, _ := userRepo.GetByEmail(ctx, "admin@mkluxe.com")
	if existing != nil {
		log.Println("Super Admin already exists, skipping seed.")
		return nil
	}

	hash, err := utils.HashPassword("admin12345")
	if err != nil {
		return err
	}

	admin := &domain.User{
		Email:        "admin@mkluxe.com",
		PasswordHash: hash,
		Name:         "System Admin",
		Role:         constants.RoleSuperAdmin,
		IsActive:     true,
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		return err
	}

	log.Println("Super Admin seeded successfully! (admin@mkluxe.com / admin12345)")
	return nil
}
