package suggestions

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/suggestion"
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/google/uuid"
)

type CreateSuggestionRequest struct {
	Body struct {
		LocationID  uuid.UUID            `json:"location_id"`
		Name        string               `json:"name" minLength:"1"`
		Category    string               `json:"category" minLength:"1"`
		Address     string               `json:"address,omitempty"`
		Phone       string               `json:"phone,omitempty"`
		WebsiteURL  string               `json:"website_url,omitempty"`
		Coordinates dto.CoordinatesInput `json:"coordinates"`
		Description string               `json:"description,omitempty"`
	}
}

type ListSuggestionsRequest struct {
	Status string `query:"status" enum:"pending,approved,rejected" default:"pending"`
}

type ApproveSuggestionRequest struct {
	ID uuid.UUID `path:"id"`
}

type RejectSuggestionRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Reason string `json:"reason" minLength:"1" maxLength:"500"`
	}
}

type SuggestionItem struct {
	ID              uuid.UUID               `json:"id"`
	UserID          uuid.UUID               `json:"user_id"`
	Nickname        string                  `json:"nickname"`
	LocationID      uuid.UUID               `json:"location_id"`
	Name            string                  `json:"name"`
	Category        string                  `json:"category"`
	Address         string                  `json:"address"`
	Phone           string                  `json:"phone"`
	WebsiteURL      string                  `json:"website_url"`
	Coordinates     coordinates.Coordinates `json:"coordinates"`
	Description     string                  `json:"description"`
	Status          string                  `json:"status"`
	RejectionReason *string                 `json:"rejection_reason,omitempty"`
	ReviewedBy      *uuid.UUID              `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time              `json:"reviewed_at,omitempty"`
	PlaceID         *uuid.UUID              `json:"place_id,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
}

type SuggestionResponse struct {
	Body SuggestionItem
}

type ListSuggestionsResponse struct {
	Body []SuggestionItem
}

func toSuggestionItem(item suggestion.Suggestion) SuggestionItem {
	return SuggestionItem{
		ID:              item.ID,
		UserID:          item.UserID,
		Nickname:        item.Nickname,
		LocationID:      item.LocationID,
		Name:            item.Name,
		Category:        item.Category,
		Address:         item.Address,
		Phone:           item.Phone,
		WebsiteURL:      item.WebsiteURL,
		Coordinates:     item.Coordinates,
		Description:     item.Description,
		Status:          item.Status,
		RejectionReason: item.RejectionReason,
		ReviewedBy:      item.ReviewedBy,
		ReviewedAt:      item.ReviewedAt,
		PlaceID:         item.PlaceID,
		CreatedAt:       item.CreatedAt,
	}
}

func toSuggestionItems(items []suggestion.Suggestion) []SuggestionItem {
	out := make([]SuggestionItem, 0, len(items))
	for _, item := range items {
		out = append(out, toSuggestionItem(item))
	}

	return out
}
