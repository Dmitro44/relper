package tasks

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/task"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error)
	Create(ctx context.Context, input CreateInput) (*task.Task, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*task.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
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
