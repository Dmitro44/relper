package tips

import (
	"github.com/Dmitro44/relper/internal/entities/tip"
	"github.com/google/uuid"
)

type CreateTipRequest struct {
	Body struct {
		LocationID uuid.UUID `json:"location_id"`
		Category   string    `json:"category" minLength:"1"`
		Title      string    `json:"title" minLength:"1"`
		Content    string    `json:"content" minLength:"1"`
		Upvotes    int32     `json:"upvotes,omitempty" minimum:"0"`
	}
}

type UpdateTipRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		LocationID *uuid.UUID `json:"location_id,omitempty"`
		Category   *string    `json:"category,omitempty" minLength:"1"`
		Title      *string    `json:"title,omitempty" minLength:"1"`
		Content    *string    `json:"content,omitempty" minLength:"1"`
		Upvotes    *int32     `json:"upvotes,omitempty" minimum:"0"`
	}
}

type GetTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type UpvoteTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type RemoveUpvoteTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type TipResponse struct {
	Body *tip.Tip
}
