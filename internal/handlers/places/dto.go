package places

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/rating"
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/google/uuid"
)

type CreatePlaceRequest struct {
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

type UpdatePlaceRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		LocationID  *uuid.UUID            `json:"location_id,omitempty"`
		Name        *string               `json:"name,omitempty" minLength:"1"`
		Category    *string               `json:"category,omitempty" minLength:"1"`
		Address     *string               `json:"address,omitempty"`
		Phone       *string               `json:"phone,omitempty"`
		WebsiteURL  *string               `json:"website_url,omitempty"`
		Coordinates *dto.CoordinatesInput `json:"coordinates,omitempty"`
		Description *string               `json:"description,omitempty"`
	}
}

type GetPlaceRequest struct {
	ID uuid.UUID `path:"id"`
}

type ListPlacesRequest struct {
	ID       uuid.UUID `path:"id"`
	Category string    `query:"category"`
	Sort     string    `query:"sort" enum:"wilson,avg_stars,votes_count,name" default:"wilson"`
}

type DeletePlaceRequest struct {
	ID uuid.UUID `path:"id"`
}

type PlaceItem struct {
	ID          uuid.UUID               `json:"id"`
	LocationID  uuid.UUID               `json:"location_id"`
	Name        string                  `json:"name"`
	Category    string                  `json:"category"`
	Address     string                  `json:"address"`
	Phone       string                  `json:"phone"`
	WebsiteURL  string                  `json:"website_url"`
	Coordinates coordinates.Coordinates `json:"coordinates"`
	Description string                  `json:"description"`
	CreatedAt   time.Time               `json:"created_at"`
}

type PlaceWithRatingItem struct {
	PlaceItem
	Rating rating.Summary `json:"rating"`
}

type PlaceResponse struct {
	Body PlaceItem
}

type PlaceWithRatingResponse struct {
	Body PlaceWithRatingItem
}

type ListPlacesResponse struct {
	Body []PlaceWithRatingItem
}

func toPlaceItem(p place.Place) PlaceItem {
	return PlaceItem{
		ID:          p.ID,
		LocationID:  p.LocationID,
		Name:        p.Name,
		Category:    p.Category,
		Address:     p.Address,
		Phone:       p.Phone,
		WebsiteURL:  p.WebsiteURL,
		Coordinates: p.Coordinates,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func toPlaceWithRatingItem(p place.WithRating) PlaceWithRatingItem {
	return PlaceWithRatingItem{
		PlaceItem: toPlaceItem(p.Place),
		Rating:    p.Rating,
	}
}

func toPlaceWithRatingItems(items []place.WithRating) []PlaceWithRatingItem {
	out := make([]PlaceWithRatingItem, 0, len(items))
	for _, item := range items {
		out = append(out, toPlaceWithRatingItem(item))
	}

	return out
}
