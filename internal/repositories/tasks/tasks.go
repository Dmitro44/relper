package tasks

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/tasks/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	Create(ctx context.Context, input CreateInput) (*model.Task, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*model.Task, error)
}

type CreateInput struct {
	LocationID  uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	PlaceID     uuid.UUID
	Done        bool
}

type UpdateInput struct {
	LocationID  *uuid.UUID
	UserID      *uuid.UUID
	Title       *string
	Description *string
	PlaceID     *uuid.UUID
	Done        *bool
}
