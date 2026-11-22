package tips

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/tip"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error)
	ListByPlace(ctx context.Context, placeID uuid.UUID) ([]tip.Tip, error)
	Create(ctx context.Context, input CreateInput) (*tip.Tip, error)
	Update(ctx context.Context, id, userID uuid.UUID, admin bool, input UpdateInput) (*tip.Tip, error)
	Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error
}

type CreateInput struct {
	LocationID *uuid.UUID
	PlaceID    *uuid.UUID
	UserID     uuid.UUID
	Category   string
	Title      string
	Content    string
}

type UpdateInput struct {
	LocationID *uuid.UUID
	PlaceID    *uuid.UUID
	Category   *string
	Title      *string
	Content    *string
}
