package audit

import (
	"context"

	auditentity "github.com/Dmitro44/relper/internal/entities/audit"
	"github.com/google/uuid"
)

type Repository interface {
	Insert(ctx context.Context, input CreateInput) error
	List(ctx context.Context, filter auditentity.Filter) ([]auditentity.Event, error)
}

type CreateInput struct {
	Source     string
	Action     string
	EntityType string
	EntityID   *uuid.UUID
	ActorID    *uuid.UUID
	Details    []byte
}
