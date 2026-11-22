package tips

import (
	"github.com/Dmitro44/relper/internal/entities/tip"
	"github.com/google/uuid"
)

type CreateTipRequest struct {
	Body struct {
		LocationID *uuid.UUID `json:"location_id,omitempty"`
		PlaceID    *uuid.UUID `json:"place_id,omitempty"`
		Category   string     `json:"category" minLength:"1"`
		Title      string     `json:"title" minLength:"1"`
		Content    string     `json:"content" minLength:"1"`
	}
}

type UpdateTipRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		LocationID *uuid.UUID `json:"location_id,omitempty"`
		PlaceID    *uuid.UUID `json:"place_id,omitempty"`
		Category   *string    `json:"category,omitempty" minLength:"1"`
		Title      *string    `json:"title,omitempty" minLength:"1"`
		Content    *string    `json:"content,omitempty" minLength:"1"`
	}
}

type GetTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type ListPlaceTipsRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type TipResponse struct {
	Body *tip.Tip
}

type TipsResponse struct {
	Body []tip.Tip
}
