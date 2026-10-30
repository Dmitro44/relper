package places

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	placesvc "github.com/Dmitro44/relper/internal/services/places"
	"github.com/danielgtaylor/huma/v2"
)

type PlaceHandler struct {
	svc *placesvc.PlaceService
}

func NewPlaceHandler(svc *placesvc.PlaceService) *PlaceHandler {
	return &PlaceHandler{
		svc: svc,
	}
}

func (h *PlaceHandler) Register(public, admin huma.API) {
	huma.Register(public, huma.Operation{
		OperationID:   "get-place",
		Method:        http.MethodGet,
		Path:          "/places/{id}",
		Summary:       "Get a place",
		Tags:          []string{"Places"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetByID)

	huma.Register(admin, huma.Operation{
		OperationID:   "update-place",
		Method:        http.MethodPatch,
		Path:          "/places/{id}",
		Summary:       "Update a place",
		Tags:          []string{"Places"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUpdate)

	huma.Register(admin, huma.Operation{
		OperationID:   "delete-place",
		Method:        http.MethodDelete,
		Path:          "/places/{id}",
		Summary:       "Delete a place",
		Tags:          []string{"Places"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)

	huma.Register(admin, huma.Operation{
		OperationID:   "create-place",
		Method:        http.MethodPost,
		Path:          "/places",
		Summary:       "Create a place",
		Tags:          []string{"Places"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)
}

func (h *PlaceHandler) HandleGetByID(ctx context.Context, in *GetPlaceRequest) (*PlaceResponse, error) {
	p, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &PlaceResponse{Body: p}, nil
}

func (h *PlaceHandler) HandleCreate(ctx context.Context, in *CreatePlaceRequest) (*PlaceResponse, error) {
	createInput := placesvc.CreateInput{
		LocationID: in.Body.LocationID,
		Name:       in.Body.Name,
		Category:   in.Body.Category,
		Address:    in.Body.Address,
		Phone:      in.Body.Phone,
		WebsiteURL: in.Body.WebsiteURL,
		Coordinates: coordinates.Coordinates{
			Lat: in.Body.Coordinates.Lat,
			Lng: in.Body.Coordinates.Lng,
		},
		Description: in.Body.Description,
		Rating:      in.Body.Rating,
	}
	p, err := h.svc.Create(ctx, createInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &PlaceResponse{Body: p}, nil
}

func (h *PlaceHandler) HandleUpdate(ctx context.Context, in *UpdatePlaceRequest) (*PlaceResponse, error) {
	updateInput := placesvc.UpdateInput{
		LocationID:  in.Body.LocationID,
		Name:        in.Body.Name,
		Category:    in.Body.Category,
		Address:     in.Body.Address,
		Phone:       in.Body.Phone,
		WebsiteURL:  in.Body.WebsiteURL,
		Description: in.Body.Description,
		Rating:      in.Body.Rating,
	}
	if in.Body.Coordinates != nil {
		updateInput.Coordinates = &coordinates.Coordinates{
			Lat: in.Body.Coordinates.Lat,
			Lng: in.Body.Coordinates.Lng,
		}
	}

	p, err := h.svc.Update(ctx, in.ID, updateInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &PlaceResponse{Body: p}, nil
}

func (h *PlaceHandler) HandleDelete(ctx context.Context, in *DeletePlaceRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
