package model

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/google/uuid"
)

type Place struct {
	ID          uuid.UUID
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Coordinates coordinates.Coordinates
	Description string
	Rating      float32
	Upvotes     int32
	CreatedAt   time.Time
}
