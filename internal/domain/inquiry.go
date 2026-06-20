package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inquiry represents a customer lead targeting a specific product
type Inquiry struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID    primitive.ObjectID `bson:"product_id" json:"product_id"`
	CustomerName string             `bson:"customer_name" json:"customer_name"`
	Phone        string             `bson:"phone" json:"phone"`
	Message      string             `bson:"message" json:"message"`
	Status       string             `bson:"status" json:"status"` // pending, reviewed, resolved
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
