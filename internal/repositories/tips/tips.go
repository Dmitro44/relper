package tips

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/tip"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error)
	Create(ctx context.Context, input CreateInput) (*tip.Tip, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*tip.Tip, error)
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
