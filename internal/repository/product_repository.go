package repository

import (
	"context"
	"time"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{collection: db.Collection("products")}
}

func (r *ProductRepository) Create(ctx context.Context, prod *domain.Product) error {
	prod.CreatedAt = time.Now()
	prod.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, prod)
	if err != nil {
		return err
	}
	prod.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var prod domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&prod)
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var prod domain.Product
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&prod)
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *ProductRepository) Update(ctx context.Context, prod *domain.Product) error {
	prod.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": prod.ID}, prod)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ProductRepository) List(ctx context.Context, filter dto.FilterRequest, page, limit int) ([]domain.Product, int64, error) {
	query := bson.M{}

	// Dynamic filtering criteria
	if filter.Status != "" {
		query["status"] = filter.Status
	}

	// Discovery Flag Filters
	if filter.IsFeatured != nil {
		query["is_featured"] = *filter.IsFeatured
	}
	if filter.IsMostSold != nil {
		query["is_most_sold"] = *filter.IsMostSold
	}

	// 💡 Category Filter: Match either primary
	if filter.CategorySlug != "" {
		query["$or"] = bson.A{
			bson.M{"category_slug": filter.CategorySlug},
		}
	}

	if filter.Search != "" {
		query["$text"] = bson.M{"$search": filter.Search}
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * limit)
	findOpts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
