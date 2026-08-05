package locations

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
	Create(ctx context.Context, input CreateInput) (*location.Location, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*location.Location, error)
}

type CreateInput struct {
	Country     string
	City        string
	Coordinates coordinates.Coordinates
}

type UpdateInput struct {
	Country     *string
	City        *string
	Coordinates *coordinates.Coordinates
}
