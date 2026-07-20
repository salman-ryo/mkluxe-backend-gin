package seed

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
	basePath := filepath.Join("internal", "seed", "products")

	seeds := []ProductSeed{
		{
			CategorySlug: "bracelets-bangles",
			ProductSlug:  "gold-tone-evil-eye-charm-bracelet",
			JSONData: `{
    "category_slug": "bracelets-bangles",
    "name": "Gold Tone Evil Eye Charm Bracelet",
    "slug": "gold-tone-evil-eye-charm-bracelet",
    "description": "Protect your energy with this stunning gold-tone bracelet. Featuring a sleek snake chain adorned with textured bead spacers, striking blue crystal teardrops, and intricate evil eye charms. Crafted with a premium anti-tarnish finish.",
    "secondary_categories": ["charms"],
    "status": "published",
    "is_featured": true,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "BRAC-EYE-001",
            "price": 649.00,
            "stock": 30,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-15-21.jpg",
            "alt_text": "Gold tone snake chain bracelet featuring evil eye charms and blue crystal teardrops",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Is the size adjustable?",
            "answer": "Yes, it features an extender chain with a secure lobster clasp for a customizable fit."
        },
        {
            "question": "Will this bracelet tarnish?",
            "answer": "No, it is treated with a high-quality anti-tarnish finish to ensure long-lasting wear and shine."
        }
    ],
    "meta_title": "Gold Tone Evil Eye Charm Bracelet | Anti-Tarnish",
    "meta_description": "Shop the Gold Tone Evil Eye Charm Bracelet. Featuring blue crystal teardrops, detailed charms, and a lasting anti-tarnish finish."
}`,
		},
		{
			CategorySlug: "necklaces",
			ProductSlug:  "gold-tone-coastal-charm-necklace",
			JSONData: `{
    "category_slug": "necklaces",
    "name": "Gold Tone Coastal Charm Necklace",
    "slug": "gold-tone-coastal-charm-necklace",
    "description": "Embrace beach-inspired elegance with this gold-tone coastal charm necklace. Featuring a delicate link chain embellished with a seashell, a faux pearl, a textured starfish, and a blue stone charm. Finished with an anti-tarnish coating for endless summer style.",
    "secondary_categories": ["charms"],
    "status": "published",
    "is_featured": false,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "NECK-CSTL-001",
            "price": 799.00,
            "stock": 25,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-15-09.jpg",
            "alt_text": "Model wearing a gold tone necklace with seashell, pearl, starfish, and blue stone charms",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Are the charms heavy?",
            "answer": "No, the charms are designed to be lightweight and comfortable for everyday wear."
        }
    ],
    "meta_title": "Gold Tone Coastal Charm Necklace | Anti-Tarnish",
    "meta_description": "Discover the Gold Tone Coastal Charm Necklace. Featuring beach-inspired shell, pearl, and starfish charms with a premium anti-tarnish finish."
}`,
		},
		{
			CategorySlug: "necklaces",
			ProductSlug:  "gold-tone-minimalist-solitaire-necklace",
			JSONData: `{
    "category_slug": "necklaces",
    "name": "Gold Tone Minimalist Solitaire Necklace",
    "slug": "gold-tone-minimalist-solitaire-necklace",
    "description": "Add a touch of subtle brilliance to your look with this gold-tone minimalist necklace. Featuring a delicate chain and a single, sparkling round-cut clear crystal in a secure setting. Crafted with a premium anti-tarnish finish for enduring everyday elegance.",
    "secondary_categories": ["pendants"],
    "status": "published",
    "is_featured": false,
    "is_most_sold": true,
    "variants": [
        {
            "sku": "NECK-SOL-002",
            "price": 499.00,
            "stock": 40,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-14-54.jpg",
            "alt_text": "Delicate gold tone chain necklace featuring a single round clear crystal pendant",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Is this necklace suitable for layering?",
            "answer": "Yes, its delicate and minimalist design makes it the perfect piece to layer with other necklaces."
        },
        {
            "question": "Will the gold tone fade quickly?",
            "answer": "No, it is protected by a durable anti-tarnish coating to maintain its radiant shine."
        }
    ],
    "meta_title": "Gold Tone Minimalist Solitaire Necklace | Anti-Tarnish",
    "meta_description": "Shop the Gold Tone Minimalist Solitaire Necklace. Featuring a delicate chain, a brilliant clear crystal pendant, and a lasting anti-tarnish finish."
}`,
		},
		{
			CategorySlug: "bracelets-bangles",
			ProductSlug:  "gold-tone-layered-heart-charm-bracelet",
			JSONData: `{
    "category_slug": "bracelets-bangles",
    "name": "Gold Tone Layered Heart Charm Bracelet",
    "slug": "gold-tone-layered-heart-charm-bracelet",
    "description": "Achieve a perfectly styled wrist stack with this gold-tone double-layered bracelet. Combining a trendy paperclip link chain with a delicate chain featuring spherical bead drops and a central polished puffed heart charm. Designed with an anti-tarnish finish.",
    "secondary_categories": [],
    "status": "published",
    "is_featured": false,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "BRAC-LAY-HRT-002",
            "price": 699.00,
            "stock": 35,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-13-40.jpg",
            "alt_text": "Model wearing a gold tone layered bracelet featuring a paperclip chain and a heart charm chain",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Is this a single piece of jewelry?",
            "answer": "Yes, this is a single bracelet designed with two integrated chains for an effortless layered look."
        },
        {
            "question": "Can I wear this everyday?",
            "answer": "Yes, it features a premium anti-tarnish finish making it highly durable for daily wear."
        }
    ],
    "meta_title": "Gold Tone Layered Heart Charm Bracelet | Anti-Tarnish",
    "meta_description": "Shop our Gold Tone Layered Heart Charm Bracelet. Featuring a trendy paperclip chain paired with a delicate heart charm chain and an anti-tarnish finish."
}`,
		},
		{
			CategorySlug: "earrings",
			ProductSlug:  "gold-tone-open-wire-heart-stud-earrings",
			JSONData: `{
    "category_slug": "earrings",
    "name": "Gold Tone Open Wire Heart Stud Earrings",
    "slug": "gold-tone-open-wire-heart-stud-earrings",
    "description": "Add a delicate touch of romance to your style with these gold-tone open wire heart stud earrings. Featuring a sleek, minimalist hollow heart silhouette and a premium anti-tarnish finish for enduring, everyday wear.",
    "secondary_categories": [],
    "status": "published",
    "is_featured": false,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "EARR-WIRE-HRT-001",
            "price": 449.00,
            "stock": 50,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-13-10.jpg",
            "alt_text": "Model wearing minimalist gold tone open wire heart stud earrings",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "Are these earrings comfortable for daily use?",
            "answer": "Yes, the open wire design makes them incredibly lightweight and perfect for all-day comfort."
        }
    ],
    "meta_title": "Gold Tone Open Wire Heart Stud Earrings | Anti-Tarnish",
    "meta_description": "Discover our Gold Tone Open Wire Heart Stud Earrings. A minimalist, lightweight design finished with a durable anti-tarnish coating for everyday style."
}`,
		},
		{
			CategorySlug: "earrings",
			ProductSlug:  "gold-tone-textured-starfish-drop-earrings",
			JSONData: `{
    "category_slug": "earrings",
    "name": "Gold Tone Textured Starfish Drop Earrings",
    "slug": "gold-tone-textured-starfish-drop-earrings",
    "description": "Channel coastal vibes with these elegant gold-tone drop earrings. Featuring detailed, textured starfish charms suspended from classic ball studs. Crafted with a high-quality anti-tarnish finish to ensure they maintain their brilliant shine.",
    "secondary_categories": [],
    "status": "published",
    "is_featured": false,
    "is_most_sold": false,
    "variants": [
        {
            "sku": "EARR-STAR-002",
            "price": 549.00,
            "stock": 30,
            "is_default": true
        }
    ],
    "media": [
        {
            "url": "photo_2026-07-19_16-13-05.jpg",
            "alt_text": "Gold tone drop earrings featuring textured starfish charms on ball studs",
            "is_primary": true
        }
    ],
    "faqs": [
        {
            "question": "What type of closure do these earrings have?",
            "answer": "They feature a secure push-back closure on a ball stud design."
        },
        {
            "question": "Will these tarnish if worn near the beach?",
            "answer": "While they feature a durable anti-tarnish finish, we recommend removing them before swimming in saltwater to prolong their lifespan."
        }
    ],
    "meta_title": "Gold Tone Textured Starfish Drop Earrings | Anti-Tarnish",
    "meta_description": "Shop the Gold Tone Textured Starfish Drop Earrings. Coastal-inspired statement pieces featuring detailed charms and a lasting anti-tarnish finish."
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
