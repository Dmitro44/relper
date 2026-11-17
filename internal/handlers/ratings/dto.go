package ratings

import (
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

type RatingResponse struct {
	Body *rating.PlaceRating
}
