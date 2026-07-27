package tags

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/tags/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	Create(ctx context.Context, input CreateInput) (*model.Tag, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*model.Tag, error)
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
