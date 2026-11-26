package place

import (
	"slices"
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/rating"
	"github.com/google/uuid"
)

var Categories = []string{"food", "transport", "healthcare", "services", "hobby", "shopping"}

func IsValidCategory(value string) bool {
	return slices.Contains(Categories, value)
}

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
