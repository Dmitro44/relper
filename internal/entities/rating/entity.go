package rating

import (
	"time"

	"github.com/google/uuid"
)

type Summary struct {
	AvgStars   float64 `json:"avg_stars"`
	VotesCount int64   `json:"votes_count"`
	Wilson     float64 `json:"wilson"`
}

type PlaceRating struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlaceID   uuid.UUID
	Stars     int16
	CreatedAt time.Time
	UpdatedAt time.Time
}
