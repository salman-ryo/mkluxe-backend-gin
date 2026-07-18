package repository

import (
	"context"
	"time"

	"mkluxe-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CategoryRepository is responsible for ALL database operations
// related to categories.
//
// It hides MongoDB implementation details from the rest of the app.
type CategoryRepository struct {

	// Reference to the MongoDB categories collection.
	//
	// Every repository method uses this collection
	// to read/write category documents.
	collection *mongo.Collection
}

// Constructor.
//
// Creates a repository that is connected to the
// "categories" collection.
func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		collection: db.Collection("categories"),
	}
}

// Create inserts a new category into MongoDB.
func (r *CategoryRepository) Create(ctx context.Context, cat *domain.Category) error {

	// Set timestamps before saving.
	// New records have identical created and updated times.
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

	// Insert the category into MongoDB.
	res, err := r.collection.InsertOne(ctx, cat)
	if err != nil {
		return err
	}

	// MongoDB automatically generates an ObjectID.
	//
	// Store that generated ID back into our struct
	// so the caller now knows the database ID.
	cat.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

// GetByID finds a category using its MongoDB ObjectID.
func (r *CategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {

	// Create an empty Category struct that MongoDB
	// will decode the document into.
	var cat domain.Category

	// Find the document whose _id matches the provided ID.
	err := r.collection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&cat)

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// GetBySlug finds a category using its slug.
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {

	var cat domain.Category

	// Query:
	//
	// {
	//     "slug": "<slug>"
	// }
	err := r.collection.
		FindOne(ctx, bson.M{"slug": slug}).
		Decode(&cat)

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// Update replaces the existing category document
// with the updated version.
func (r *CategoryRepository) Update(ctx context.Context, cat *domain.Category) error {

	// Since we're modifying the category,
	// update the last-modified timestamp.
	cat.UpdatedAt = time.Now()

	// Replace the document whose _id matches cat.ID
	// with the contents of cat.
	//
	// ReplaceOne replaces the entire document.
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": cat.ID},
		cat,
	)

	return err
}

// Delete removes a category by its ObjectID.
func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {

	// Delete the document whose _id matches.
	_, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	return err
}

// ListAll returns every category sorted by sort_order.
// ListAll returns categories sorted by sort_order. Supports optional filtering by isFeatured.
func (r *CategoryRepository) ListAll(ctx context.Context, isFeatured *bool) ([]domain.Category, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "sort_order", Value: 1}})

	filter := bson.M{}

	// 💡 If a filter is provided, add it to the MongoDB query
	if isFeatured != nil {
		filter["is_featured"] = *isFeatured
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
