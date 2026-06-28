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

type InquiryRepository struct {
	collection *mongo.Collection
}

func NewInquiryRepository(db *mongo.Database) *InquiryRepository {
	return &InquiryRepository{collection: db.Collection("inquiries")}
}

func (r *InquiryRepository) Create(ctx context.Context, inq *domain.Inquiry) error {
	inq.CreatedAt = time.Now()
	inq.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, inq)
	if err != nil {
		return err
	}
	inq.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *InquiryRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *InquiryRepository) List(ctx context.Context, page, limit int) ([]domain.Inquiry, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * limit)
	findOpts := options.Find().SetSkip(skip).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var inquiries []domain.Inquiry
	if err := cursor.All(ctx, &inquiries); err != nil {
		return nil, 0, err
	}
	return inquiries, total, nil
}
