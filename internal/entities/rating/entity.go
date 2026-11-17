package rating

import (
	"time"

	"github.com/google/uuid"
)

type PlaceRating struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlaceID   uuid.UUID
	Stars     int16
	CreatedAt time.Time
	UpdatedAt time.Time
}
