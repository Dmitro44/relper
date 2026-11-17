package like

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TipID     *uuid.UUID
	CommentID *uuid.UUID
	CreatedAt time.Time
}

type LikeSummary struct {
	Count  int64 `json:"count"`
	MyLike *bool `json:"my_like,omitempty"`
}
