package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsRepository struct {
	db *mongo.Database
}

func NewStatsRepository(db *mongo.Database) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) CountProducts(ctx context.Context) (int64, error) {
	return r.db.Collection("products").CountDocuments(ctx, bson.M{})
}

func (r *StatsRepository) CountCategories(ctx context.Context) (int64, error) {
	return r.db.Collection("categories").CountDocuments(ctx, bson.M{})
}

func (r *StatsRepository) CountInquiries(ctx context.Context) (int64, error) {
	return r.db.Collection("inquiries").CountDocuments(ctx, bson.M{})
}

func (r *StatsRepository) CountFeaturedProducts(ctx context.Context) (int64, error) {
	return r.db.Collection("products").CountDocuments(ctx, bson.M{"is_featured": true})
}

func (r *StatsRepository) CountMostSoldProducts(ctx context.Context) (int64, error) {
	return r.db.Collection("products").CountDocuments(ctx, bson.M{"is_most_sold": true})
}

func (r *StatsRepository) GetProductStatusCounts(ctx context.Context) (map[string]int64, error) {
	col := r.db.Collection("products")
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	counts := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err == nil {
			status := result.ID
			if status == "" {
				status = "unknown"
			}
			counts[status] = result.Count
		}
	}
	return counts, nil
}

func (r *StatsRepository) GetInquiryStatusCounts(ctx context.Context) (map[string]int64, error) {
	col := r.db.Collection("inquiries")
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	counts := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err == nil {
			status := result.ID
			if status == "" {
				status = "unknown"
			}
			counts[status] = result.Count
		}
	}
	return counts, nil
}

func (r *StatsRepository) GetProductCategoryCounts(ctx context.Context) (map[string]int64, error) {
	col := r.db.Collection("products")
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category_slug"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	counts := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err == nil {
			slug := result.ID
			if slug != "" {
				counts[slug] = result.Count
			}
		}
	}
	return counts, nil
}

func (r *StatsRepository) GetStockStats(ctx context.Context) (totalStock int64, outOfStockCount int64, err error) {
	col := r.db.Collection("products")

	// Aggregate total stock
	stockPipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$variants"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_stock", Value: bson.D{{Key: "$sum", Value: "$variants.stock"}}},
		}}},
	}

	stockCursor, err := col.Aggregate(ctx, stockPipeline)
	if err == nil {
		defer stockCursor.Close(ctx)
		if stockCursor.Next(ctx) {
			var res struct {
				TotalStock int64 `bson:"total_stock"`
			}
			if err := stockCursor.Decode(&res); err == nil {
				totalStock = res.TotalStock
			}
		}
	}

	// Aggregate products that are completely out of stock (sum of all variants' stock is 0)
	outOfStockPipeline := mongo.Pipeline{
		// Group by product ID, summing variant stocks
		{{Key: "$project", Value: bson.D{
			{Key: "total_product_stock", Value: bson.D{{Key: "$sum", Value: "$variants.stock"}}},
		}}},
		{{Key: "$match", Value: bson.M{"total_product_stock": 0}}},
		{{Key: "$count", Value: "count"}},
	}

	outCursor, err := col.Aggregate(ctx, outOfStockPipeline)
	if err == nil {
		defer outCursor.Close(ctx)
		if outCursor.Next(ctx) {
			var res struct {
				Count int64 `bson:"count"`
			}
			if err := outCursor.Decode(&res); err == nil {
				outOfStockCount = res.Count
			}
		}
	}

	return totalStock, outOfStockCount, nil
}
