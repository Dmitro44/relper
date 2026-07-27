package model

import (
	"time"

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
	Coordinates Coordinates
	Description string
	Rating      float32
	Upvotes     int32
	CreatedAt   time.Time
}

type Coordinates struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}
