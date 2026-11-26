package suggestion

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/google/uuid"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type Suggestion struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Nickname    string
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Description string
	coordinates.Coordinates
	Status          string
	RejectionReason *string
	ReviewedBy      *uuid.UUID
	ReviewedAt      *time.Time
	PlaceID         *uuid.UUID
	CreatedAt       time.Time
}
