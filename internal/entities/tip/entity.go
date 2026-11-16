package tip

import (
	"time"

	"github.com/google/uuid"
)

type Tip struct {
	ID         uuid.UUID
	LocationID *uuid.UUID
	PlaceID    *uuid.UUID
	UserID     uuid.UUID
	Category   string
	Title      string
	Content    string
	CreatedAt  time.Time
}
