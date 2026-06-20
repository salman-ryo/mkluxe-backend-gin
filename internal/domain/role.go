package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role holds metadata for permissions if stored in DB, otherwise managed via constants
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Permissions []string           `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
