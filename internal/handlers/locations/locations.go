package locations

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	locationsvc "github.com/Dmitro44/relper/internal/services/locations"
	"github.com/danielgtaylor/huma/v2"
)

type LocationHandler struct {
	svc *locationsvc.LocationService
}

func NewLocationHandler(svc *locationsvc.LocationService) *LocationHandler {
	return &LocationHandler{
		svc: svc,
	}
}

func (h *LocationHandler) Register(public, admin huma.API) {
	huma.Register(public, huma.Operation{
		OperationID:   "get-location",
		Method:        http.MethodGet,
		Path:          "/locations/{id}",
		Summary:       "Get a location",
		Tags:          []string{"Locations"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetByID)

	huma.Register(admin, huma.Operation{
		OperationID:   "create-location",
		Method:        http.MethodPost,
		Path:          "/locations",
		Summary:       "Create a location",
		Tags:          []string{"Locations"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)

	huma.Register(admin, huma.Operation{
		OperationID:   "update-location",
		Method:        http.MethodPatch,
		Path:          "/locations/{id}",
		Summary:       "Update a location",
		Tags:          []string{"Locations"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUpdate)

	huma.Register(admin, huma.Operation{
		OperationID:   "delete-location",
		Method:        http.MethodDelete,
		Path:          "/locations/{id}",
		Summary:       "Delete a location",
		Tags:          []string{"Locations"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)
}

func (h *LocationHandler) HandleGetByID(ctx context.Context, in *GetLocationRequest) (*LocationResponse, error) {
	loc, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &LocationResponse{Body: loc}, nil
}

func (h *LocationHandler) HandleCreate(ctx context.Context, in *CreateLocationRequest) (*LocationResponse, error) {
	createInput := locationsvc.CreateInput{
		Country: in.Body.Country,
		City:    in.Body.City,
		Coordinates: coordinates.Coordinates{
			Lat: in.Body.Coordinates.Lat,
			Lng: in.Body.Coordinates.Lng,
		},
	}
	loc, err := h.svc.Create(ctx, createInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &LocationResponse{Body: loc}, nil
}

func (h *LocationHandler) HandleUpdate(ctx context.Context, in *UpdateLocationRequest) (*LocationResponse, error) {
	updateInput := locationsvc.UpdateInput{
		Country: in.Body.Country,
		City:    in.Body.City,
	}
	if in.Body.Coordinates != nil {
		updateInput.Coordinates = &coordinates.Coordinates{
			Lat: in.Body.Coordinates.Lat,
			Lng: in.Body.Coordinates.Lng,
		}
	}

	loc, err := h.svc.Update(ctx, in.ID, updateInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &LocationResponse{Body: loc}, nil
}

func (h *LocationHandler) HandleDelete(ctx context.Context, in *DeleteLocationRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
