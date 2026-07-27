package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	LocationID  uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	PlaceID     uuid.UUID
	Done        bool
	CreatedAt   time.Time
}
