package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ProductSeed holds the pathing info and the raw JSON payload
type ProductSeed struct {
	CategorySlug string
	ProductSlug  string
	JSONData     string
}

func main() {
	// Base path provided in your prompt
	basePath := filepath.Join("C:", "Users", "salma", "Development", "Jiyu", "Golang", "Projects", "mkluxe-backend", "internal", "seed", "products")

	seeds := []ProductSeed{
		{
			CategorySlug: "earrings",
			ProductSlug:  "gold-tone-textured-square-hoop-earrings",
			JSONData: `{
    "category_slug": "earrings",
    "name": "Gold Tone Textured Square Hoop Earrings",
    "slug": "gold-tone-textured-square-hoop-earrings",
    "description": "Make a sophisticated statement with these gold-tone square hoop earrings. Featuring a unique geometric shape with a hammered texture and a premium anti-tarnish finish for an elegant, everyday look.",
    "secondary_categories": [],
    "status": "published",
    "is_featured": false,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "EARR-SQR-HOOP-001",
            "price": 599.00,
            "stock": 25,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-05-51.jpg",
            "alt_text": "Model wearing gold tone textured square geometric hoop earrings",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Are these hoops heavy?",
            "answer": "Despite their bold geometric design, these earrings are crafted to be comfortably lightweight for all-day wear."
        },
        {
            "question": "Will the gold finish fade over time?",
            "answer": "No, they are treated with a durable anti-tarnish finish to ensure they maintain their brilliant shine."
        }
    ],
    "meta_title": "Gold Tone Textured Square Hoop Earrings | Anti-Tarnish",
    "meta_description": "Shop the Gold Tone Textured Square Hoop Earrings. Featuring a vintage-inspired hammered texture, geometric design, and a lasting anti-tarnish finish."
}`,
		},
		{
			CategorySlug: "earrings",
			ProductSlug:  "gold-tone-melting-heart-statement-earrings",
			JSONData: `{
    "category_slug": "earrings",
    "name": "Gold Tone Melting Heart Statement Earrings",
    "slug": "gold-tone-melting-heart-statement-earrings",
    "description": "Embrace avant-garde style with these gold-tone melting heart earrings. Featuring a high-polish, contemporary drip design, these statement studs are finished with an anti-tarnish coating for lasting brilliance and modern flair.",
    "secondary_categories": [],
    "status": "published",
    "is_featured": true,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "EARR-MELT-HRT-001",
            "price": 699.00,
            "stock": 20,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-05-45.jpg",
            "alt_text": "Model wearing bold gold tone melting heart drip statement earrings",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Are these statement earrings comfortable?",
            "answer": "Yes, they are designed as stud earrings to provide a secure and comfortable fit while offering a bold, elongated look."
        }
    ],
    "meta_title": "Gold Tone Melting Heart Statement Earrings | Anti-Tarnish",
    "meta_description": "Discover our Gold Tone Melting Heart Statement Earrings. A bold, contemporary liquid metal design finished with a premium anti-tarnish coating."
}`,
		},
		{
			CategorySlug: "bracelets-bangles",
			ProductSlug:  "gold-tone-layered-snake-chain-beaded-bracelet",
			JSONData: `{
    "category_slug": "bracelets-bangles",
    "name": "Gold Tone Layered Snake Chain Beaded Bracelet",
    "slug": "gold-tone-layered-snake-chain-beaded-bracelet",
    "description": "Elevate your wrist stack effortlessly with this gold-tone layered bracelet. Showcasing a flat herringbone snake chain and a rounded snake chain accented with polished sphere beads. Crafted with a premium anti-tarnish finish for enduring, everyday elegance.",
    "secondary_categories": [],
    "status": "published",
    "is_featured": false,
    "is_most_sold": true,
    "variants": [
        {
            "sku": "BRAC-LAY-SNK-001",
            "price": 749.00,
            "stock": 35,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-05-37.jpg",
            "alt_text": "Model wearing a gold tone layered snake chain and herringbone bracelet with spherical bead accents",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Are these two separate bracelets?",
            "answer": "No, this is a beautifully designed single bracelet featuring two attached chains for a perfect, tangle-free layered look."
        },
        {
            "question": "Is this bracelet anti-tarnish?",
            "answer": "Yes, it is protected with a high-quality anti-tarnish finish to preserve its radiant golden shine."
        }
    ],
    "meta_title": "Gold Tone Layered Snake Chain Bracelet | Anti-Tarnish",
    "meta_description": "Shop the Gold Tone Layered Snake Chain Beaded Bracelet. Featuring a chic flat herringbone and rounded chain design with a lasting anti-tarnish finish."
}`,
		},
	}

	for _, seed := range seeds {
		// 1. Construct the folder path
		targetDir := filepath.Join(basePath, seed.CategorySlug, seed.ProductSlug)

		// 2. Create directories (equivalent to mkdir -p)
		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", targetDir, err)
		}

		// 3. Construct the file path
		filePath := filepath.Join(targetDir, "data.json")

		// 4. Write the JSON payload into data.json
		err = os.WriteFile(filePath, []byte(seed.JSONData), 0644)
		if err != nil {
			log.Fatalf("Failed to write file %s: %v", filePath, err)
		}

		fmt.Printf("Successfully generated: %s\n", filePath)
	}

	fmt.Println("---")
	fmt.Println("All product seed files generated successfully!")
}
