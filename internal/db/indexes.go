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

	// Clean up the timeout context when the function exits.
	//
	// This stops the timer and releases resources associated
	// with the context.
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

		// Unique index on slug.
		//
		// Example:
		// /products/iphone-15
		//
		// Every product must have a unique slug.
		{
			Keys: bson.D{
				{Key: "slug", Value: 1},
			},
			Options: uniqueOpt,
		},

		// Compound index:
		// status + primary_category_id
		//
		// Useful for queries like:
		//
		// db.products.find({
		//     status: "active",
		//     primary_category_id: "electronics"
		// })
		//
		// MongoDB can use this index instead of scanning
		// every product document.
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "primary_category_id", Value: 1},
			},
		},

		// Text index for full-text search.
		//
		// Allows queries like:
		//
		// db.products.find({
		//     $text: { $search: "wireless headphones" }
		// })
		//
		// MongoDB will search both name and description.
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
			},
		},
	})

	// If index creation failed, stop immediately.
	if err != nil {
		return err
	}

	// ============================================================
	// 2. CATEGORY INDEXES
	// ============================================================

	// Get the categories collection.
	categoryCol := db.Collection("categories")

	// Create a unique index on slug.
	//
	// Example:
	//
	// electronics
	// phones
	// laptops
	//
	// No duplicate category slugs allowed.
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

	// Create a unique index on email.
	//
	// This prevents duplicate accounts:
	//
	// john@example.com
	// john@example.com  <- rejected
	//
	// MongoDB itself enforces this rule.
	_, err = userCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: uniqueOpt,
	})

	if err != nil {
		return err
	}

	// If we reach this point, all indexes were created
	// successfully (or already existed).
	log.Println("MongoDB indexes verified successfully.")

	// nil means no error occurred.
	return nil
}
