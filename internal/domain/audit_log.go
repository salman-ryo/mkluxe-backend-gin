package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog tracks administrative actions
type AuditLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActorID    primitive.ObjectID `bson:"actor_id" json:"actor_id"`
	Action     string             `bson:"action" json:"action"`
	EntityType string             `bson:"entity_type" json:"entity_type"` // e.g., "product", "user"
	EntityID   primitive.ObjectID `bson:"entity_id" json:"entity_id"`
	Changes    interface{}        `bson:"changes,omitempty" json:"changes,omitempty"` // Flexible BSON object
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
