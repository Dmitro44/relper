package places

import (
	"github.com/Dmitro44/relper/internal/entities/place"
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
		Rating      float32              `json:"rating,omitempty" minimum:"0" maximum:"5"`
		Upvotes     int32                `json:"upvotes,omitempty" minimum:"0"`
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
		Rating      *float32              `json:"rating,omitempty" minimum:"0" maximum:"5"`
		Upvotes     *int32                `json:"upvotes,omitempty" minimum:"0"`
	}
}

type GetPlaceRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeletePlaceRequest struct {
	ID uuid.UUID `path:"id"`
}

type PlaceResponse struct {
	Body *place.Place
}
