package places

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/google/uuid"
)

type RatingCounts struct {
	AvgStars     float64
	VotesCount   int64
	Recommend    int64
	NotRecommend int64
}

type PlaceWithRating struct {
	place.Place
	RatingCounts
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error)
	GetByIDWithRating(ctx context.Context, id uuid.UUID) (*PlaceWithRating, error)
	ListByLocation(ctx context.Context, locationID uuid.UUID, category string) ([]PlaceWithRating, error)
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
}
