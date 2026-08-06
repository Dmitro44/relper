package locations

import (
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
	Body *location.Location
}
