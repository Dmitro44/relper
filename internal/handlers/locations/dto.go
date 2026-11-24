package locations

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/google/uuid"
)

type GetLocationRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteLocationRequest struct {
	ID uuid.UUID `path:"id"`
}

type CreateLocationRequest struct {
	Body struct {
		Country     string               `json:"country" minLength:"1"`
		City        string               `json:"city" minLength:"1"`
		Coordinates dto.CoordinatesInput `json:"coordinates"`
	}
}

type UpdateLocationRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Country     *string               `json:"country,omitempty" minLength:"1"`
		City        *string               `json:"city,omitempty" minLength:"1"`
		Coordinates *dto.CoordinatesInput `json:"coordinates,omitempty"`
	}
}

type LocationResponse struct {
	Body LocationItem
}

type ListLocationsRequest struct{}

type LocationItem struct {
	ID          uuid.UUID               `json:"id"`
	Country     string                  `json:"country"`
	City        string                  `json:"city"`
	Coordinates coordinates.Coordinates `json:"coordinates"`
	CreatedAt   time.Time               `json:"created_at"`
}

type LocationsResponse struct {
	Body []LocationItem
}

func toLocationItem(item location.Location) LocationItem {
	return LocationItem{
		ID:          item.ID,
		Country:     item.Country,
		City:        item.City,
		Coordinates: item.Coordinates,
		CreatedAt:   item.CreatedAt,
	}
}

func toLocationItems(items []location.Location) []LocationItem {
	out := make([]LocationItem, 0, len(items))
	for _, item := range items {
		out = append(out, toLocationItem(item))
	}

	return out
}
