package ratings

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/rating"
	"github.com/google/uuid"
)

type Repository interface {
	Upsert(ctx context.Context, userID, placeID uuid.UUID, stars int16) (*rating.PlaceRating, error)
	Delete(ctx context.Context, userID, placeID uuid.UUID) error
}
