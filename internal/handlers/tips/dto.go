package tips

import (
	"time"

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

type TipItem struct {
	ID         uuid.UUID  `json:"id"`
	LocationID *uuid.UUID `json:"location_id,omitempty"`
	PlaceID    *uuid.UUID `json:"place_id,omitempty"`
	UserID     uuid.UUID  `json:"user_id"`
	Category   string     `json:"category"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
}

type TipResponse struct {
	Body TipItem
}

type TipsResponse struct {
	Body []TipItem
}

func toTipItem(t tip.Tip) TipItem {
	return TipItem{
		ID:         t.ID,
		LocationID: t.LocationID,
		PlaceID:    t.PlaceID,
		UserID:     t.UserID,
		Category:   t.Category,
		Title:      t.Title,
		Content:    t.Content,
		CreatedAt:  t.CreatedAt,
	}
}

func toTipItems(items []tip.Tip) []TipItem {
	out := make([]TipItem, 0, len(items))
	for _, item := range items {
		out = append(out, toTipItem(item))
	}

	return out
}
