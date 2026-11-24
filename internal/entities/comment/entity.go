package comment

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID
	PlaceID   *uuid.UUID
	TipID     *uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
}

type WithAuthor struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Nickname  string
	Content   string
	CreatedAt time.Time
}
