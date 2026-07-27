package tips

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/tips/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Tip, error)
	Create(ctx context.Context, input CreateInput) (*model.Tip, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*model.Tip, error)
}

type CreateInput struct {
	LocationID uuid.UUID
	UserID     uuid.UUID
	Category   string
	Title      string
	Content    string
	Upvotes    int32
}

type UpdateInput struct {
	LocationID *uuid.UUID
	UserID     *uuid.UUID
	Category   *string
	Title      *string
	Content    *string
	Upvotes    *int32
}
