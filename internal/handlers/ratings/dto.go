package ratings

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/rating"
	"github.com/google/uuid"
)

type SetRatingRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Stars int16 `json:"stars" minimum:"1" maximum:"5"`
	}
}

type DeleteRatingRequest struct {
	ID uuid.UUID `path:"id"`
}

type RatingItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	PlaceID   uuid.UUID `json:"place_id"`
	Stars     int16     `json:"stars"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RatingResponse struct {
	Body RatingItem
}

func toRatingItem(r rating.PlaceRating) RatingItem {
	return RatingItem{
		ID:        r.ID,
		UserID:    r.UserID,
		PlaceID:   r.PlaceID,
		Stars:     r.Stars,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
