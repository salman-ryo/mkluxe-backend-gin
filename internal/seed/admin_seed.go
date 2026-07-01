package seed

import (
	"context"
	"errors"
	"log"
	"os" // 👈 Added to read environment variables

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

	// 1. Load admin credentials from environment variables with fallbacks
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@mkluxe.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin12345"
	}

	// 2. Fetch the user checking the dynamic email
	existing, err := userRepo.GetByEmail(ctx, adminEmail)

	// If there's an error, check if it's just telling us the user doesn't exist.
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	// 3. If the user was actually found, skip the seeding process safely
	if existing != nil {
		log.Println("Super Admin already exists, skipping seed.")
		return nil
	}

	// 4. Admin doesn't exist, hash the dynamic password and proceed to create it
	hash, err := utils.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	admin := &domain.User{
		Email:        adminEmail,
		PasswordHash: hash,
		Name:         "System Admin",
		Role:         constants.RoleSuperAdmin,
		IsActive:     true,
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		return err
	}

	log.Printf("Super Admin seeded successfully! (%s / %s)", adminEmail, adminPassword)
	return nil
}
