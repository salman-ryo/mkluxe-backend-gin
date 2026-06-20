package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents an admin panel account
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"` // "-" prevents accidental JSON exposure
	Name         string             `bson:"name" json:"name"`
	Role         string             `bson:"role" json:"role"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	LastLogin    time.Time          `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
