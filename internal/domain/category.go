package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category is a flat taxonomy representation
type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Slug        string             `bson:"slug" json:"slug"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	ImageUrl    string             `bson:"image_url" json:"image_url"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	IsFeatured  bool               `bson:"is_featured" json:"is_featured"`
	SortOrder   int                `bson:"sort_order" json:"sort_order"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
