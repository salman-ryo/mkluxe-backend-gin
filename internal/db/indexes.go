package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates all indexes required by the application.
//
// Why?
// - Improves query performance.
// - Enforces uniqueness constraints.
// - Ensures indexes exist when the app starts.
//
// This function is safe to run repeatedly because MongoDB
// won't recreate identical indexes that already exist.
func EnsureIndexes(db *mongo.Database) error {

	// Create a context with a 15-second timeout.
	//
	// Index creation can take some time, especially on larger collections.
	// We don't want this operation hanging forever if MongoDB becomes
	// unreachable or something goes wrong.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create reusable index options.
	//
	// SetUnique(true) means MongoDB will reject duplicate values
	// for any index that uses these options.
	//
	// Example:
	// If slug = "iphone-15" already exists,
	// inserting another document with the same slug will fail.
	uniqueOpt := options.Index().SetUnique(true)

	// ============================================================
	// 1. PRODUCT INDEXES
	// ============================================================

	// Get a reference to the "products" collection.
	productCol := db.Collection("products")

	// Create multiple indexes in a single operation.
	_, err := productCol.Indexes().CreateMany(ctx, []mongo.IndexModel{

		// --------------------------------------------------------
		// Unique product slug
		// --------------------------------------------------------
		//
		// Used for product URLs such as:
		//
		// /products/iphone-15
		//
		// Every product must have a unique slug.
		{
			Keys: bson.D{
				{Key: "slug", Value: 1},
			},
			Options: uniqueOpt,
		},

		// --------------------------------------------------------
		// Compound index:
		// status + category_slug
		// --------------------------------------------------------
		//
		// Optimizes queries such as:
		//
		// db.products.find({
		//     status: "active",
		//     category_slug: "electronics"
		// })
		//
		// This is useful for category listing pages where only
		// active products are displayed.
		//
		// In production, products reference the category by its
		// unique slug rather than a category ID.
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "category_slug", Value: 1},
			},
		},

		// --------------------------------------------------------
		// Text search index
		// --------------------------------------------------------
		//
		// Enables MongoDB full-text search.
		//
		// Example:
		//
		// db.products.find({
		//     $text: { $search: "wireless headphones" }
		// })
		//
		// MongoDB searches both the product name and description.
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
			},
		},
	})

	// Stop immediately if product index creation failed.
	if err != nil {
		return err
	}

	// ============================================================
	// 2. CATEGORY INDEXES
	// ============================================================

	// Get the categories collection.
	categoryCol := db.Collection("categories")

	// Unique category slug.
	//
	// Examples:
	//
	// electronics
	// laptops
	// smartphones
	//
	// Products reference categories using this slug, so it must
	// remain unique across the collection.
	_, err = categoryCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "slug", Value: 1},
		},
		Options: uniqueOpt,
	})

	if err != nil {
		return err
	}

	// ============================================================
	// 3. USER INDEXES
	// ============================================================

	// Get the users collection.
	userCol := db.Collection("users")

	// Unique email address.
	//
	// Prevents duplicate user accounts:
	//
	// john@example.com
	// john@example.com  <- rejected
	//
	// MongoDB enforces this uniqueness automatically.
	_, err = userCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: uniqueOpt,
	})

	if err != nil {
		return err
	}

	// All indexes now exist (either newly created or already present).
	log.Println("MongoDB indexes verified successfully.")

	return nil
}
