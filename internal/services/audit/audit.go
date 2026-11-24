package audit

import (
	"context"
	"encoding/json"

	auditentity "github.com/Dmitro44/relper/internal/entities/audit"
	apperr "github.com/Dmitro44/relper/internal/errors"
	auditrepo "github.com/Dmitro44/relper/internal/repositories/audit"
)

type AuditService struct {
	repo auditrepo.Repository
}

func NewAuditService(repo auditrepo.Repository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Record(ctx context.Context, input auditentity.EventInput) error {
	details := []byte("{}")
	if input.Details != nil {
		encoded, err := json.Marshal(input.Details)
		if err != nil {
			return apperr.NewInternalError("failed to encode audit details", err)
		}
		details = encoded
	}

	err := s.repo.Insert(ctx, auditrepo.CreateInput{
		Source:     "app",
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		ActorID:    input.ActorID,
		Details:    details,
	})
	if err != nil {
		return apperr.NewInternalError("failed to record audit event", err)
	}

	return nil
}

func (s *AuditService) List(ctx context.Context, filter auditentity.Filter) ([]auditentity.Event, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}

	items, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, apperr.NewInternalError("failed to list audit events", err)
	}

	return items, nil
}
