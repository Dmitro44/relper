package tip

import (
	"time"

	"github.com/google/uuid"
)

type Tip struct {
	ID         uuid.UUID
	LocationID uuid.UUID
	UserID     uuid.UUID
	Category   string
	Title      string
	Content    string
	Upvotes    int32
	CreatedAt  time.Time
}
