package seed

import (
	"context"
	"log"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

// SeedCategories populates the database with the baseline jewelry categories with SEO-optimized descriptions
func SeedCategories(db *mongo.Database) error {
	categoryRepo := repository.NewCategoryRepository(db)
	ctx := context.Background()

	// Define the base categories with enhanced, SEO-friendly descriptions
	categories := []struct {
		Name        string
		Description string
	}{
		{
			Name:        "Earrings",
			Description: "Explore our exquisite collection of earrings, featuring everyday elegant studs, classic hoops, traditional Indian jhumkas, and statement danglers for every occasion.",
		},
		{
			Name:        "Necklaces",
			Description: "Discover stunning necklaces to elevate your style, ranging from delicate everyday chains to bold, intricately designed chokers perfect for ethnic and modern wear.",
		},
		{
			Name:        "Jewelry Sets",
			Description: "Shop perfectly matched jewelry sets, including beautiful necklace and earring combos designed for effortless elegance, parties, and special events.",
		},
		{
			Name:        "Bracelets & Bangles",
			Description: "Adorn your wrists with our premium collection of wristwear, featuring delicate bracelets, traditional kadas, and beautifully crafted bangles.",
		},
		{
			Name:        "Anklets",
			Description: "Step out in style with our elegant anklets and intricate toe chains, blending traditional Indian charm with modern grace.",
		},
		{
			Name:        "Nose Pins",
			Description: "Accentuate your look with our diverse range of nose jewelry, featuring subtle, elegant nose studs and striking septum rings.",
		},
		{
			Name:        "Hair Accessories",
			Description: "Complete your festive and bridal ensembles with our stunning hair accessories, including traditional maang tikkas and decorative hair pins.",
		},
		{
			Name:        "Bridal Jewelry",
			Description: "Make your special day unforgettable with our luxurious bridal jewelry collections and complete wedding sets crafted for the elegant bride.",
		},
		{
			Name:        "Traditional Jewelry",
			Description: "Celebrate heritage with our authentic traditional jewelry, showcasing exquisite South Indian temple designs and timeless ethnic pieces.",
		},
		{
			Name:        "Men's Jewelry",
			Description: "Browse our exclusive men's jewelry collection, featuring bold chains, sophisticated bracelets, and classic rings designed for the modern man.",
		},
		{
			Name:        "Pendants",
			Description: "Find the perfect centerpiece with our versatile selection of designer pendants and pendant-only pieces, ideal for pairing with your favorite chains.",
		},
		{
			Name:        "Mangalsutra",
			Description: "Honor your sacred bond with our beautifully crafted mangalsutras, blending traditional significance with contemporary, everyday wearable designs.",
		},
	}

	for _, catData := range categories {
		slug := utils.GenerateSlug(catData.Name)

		// Check if category already exists to prevent duplicates on multiple runs
		existing, _ := categoryRepo.GetBySlug(ctx, slug)
		if existing != nil {
			log.Printf("Category '%s' already exists, skipping.", catData.Name)
			continue
		}

		newCategory := &domain.Category{
			Name:        catData.Name,
			Slug:        slug,
			IsActive:    true,
			Description: catData.Description,
		}

		if err := categoryRepo.Create(ctx, newCategory); err != nil {
			log.Printf("Error creating category '%s': %v", catData.Name, err)
			return err
		}

		log.Printf("Successfully seeded category: %s", catData.Name)
	}

	log.Println("All baseline categories seeded successfully!")
	return nil
}
