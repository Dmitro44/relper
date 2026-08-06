package places

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error)
	Create(ctx context.Context, input CreateInput) (*place.Place, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*place.Place, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Coordinates coordinates.Coordinates
	Description string
	Rating      float32
}

type UpdateInput struct {
	LocationID  *uuid.UUID
	Name        *string
	Category    *string
	Address     *string
	Phone       *string
	WebsiteURL  *string
	Coordinates *coordinates.Coordinates
	Description *string
	Rating      *float32
}
