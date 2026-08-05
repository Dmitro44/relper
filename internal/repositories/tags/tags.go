package tags

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/tag"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*tag.Tag, error)
	Create(ctx context.Context, input CreateInput) (*tag.Tag, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*tag.Tag, error)
	AddToPlace(ctx context.Context, placeID, tagID, userID uuid.UUID) error
	RemoveFromPlace(ctx context.Context, placeID, tagID uuid.UUID) error
}

type CreateInput struct {
	Name string
	Slug string
}

type UpdateInput struct {
	Name *string
	Slug *string
}
