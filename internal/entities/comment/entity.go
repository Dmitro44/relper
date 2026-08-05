package comment

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID
	TargetType string
	TargetID   uuid.UUID
	UserID     uuid.UUID
	Content    string
	CreatedAt  time.Time
}
