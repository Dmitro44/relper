package places

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/places/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Place, error)
	Create(ctx context.Context, input CreateInput) (*model.Place, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*model.Place, error)
}

type CreateInput struct {
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Coordinates model.Coordinates
	Description string
	Rating      float32
	Upvotes     int32
}

type UpdateInput struct {
	LocationID  *uuid.UUID
	Name        *string
	Category    *string
	Address     *string
	Phone       *string
	WebsiteURL  *string
	Coordinates *model.Coordinates
	Description *string
	Rating      *float32
	Upvotes     *int32
}
