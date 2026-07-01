package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product is the primary aggregate root for the jewelry catalog
type Product struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Slug                string               `bson:"slug" json:"slug"`
	Name                string               `bson:"name" json:"name"`
	Description         string               `bson:"description" json:"description"`
	PrimaryCategoryID   primitive.ObjectID   `bson:"primary_category_id" json:"primary_category_id"`
	SecondaryCategories []primitive.ObjectID `bson:"secondary_categories,omitempty" json:"secondary_categories,omitempty"`

	Status string `bson:"status" json:"status"` // draft, published, archived

	// New Discovery Flags
	IsFeatured bool `bson:"is_featured" json:"is_featured"`
	IsMostSold bool `bson:"is_most_sold" json:"is_most_sold"`

	Variants []Variant `bson:"variants" json:"variants"`
	Media    []Media   `bson:"media" json:"media"`
	FAQs     []FAQ     `bson:"faqs,omitempty" json:"faqs,omitempty"`

	// SEO Fields
	MetaTitle       string `bson:"meta_title,omitempty" json:"meta_title,omitempty"`
	MetaDescription string `bson:"meta_description,omitempty" json:"meta_description,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
