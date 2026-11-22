package place

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/rating"
	"github.com/google/uuid"
)

type Place struct {
	ID                      uuid.UUID
	LocationID              uuid.UUID
	Name                    string
	Category                string
	Address                 string
	Phone                   string
	WebsiteURL              string
	coordinates.Coordinates `json:"Coordinates"`
	Description             string
	CreatedAt               time.Time
}

type WithRating struct {
	Place
	Rating rating.Summary `json:"rating"`
}
