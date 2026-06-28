package service

import (
	"context"

	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditService struct {
	repo *repository.AuditLogRepository
}

func NewAuditService(repo *repository.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

// LogAction captures and saves an event to the audit trail
func (s *AuditService) LogAction(ctx context.Context, actorID primitive.ObjectID, action, entityType string, entityID primitive.ObjectID, changes interface{}) error {
	log := &domain.AuditLog{
		ActorID:    actorID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Changes:    changes,
	}
	return s.repo.Create(ctx, log)
}
