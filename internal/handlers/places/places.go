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

	huma.Register(public, huma.Operation{
		OperationID:   "list-location-places",
		Method:        http.MethodGet,
		Path:          "/locations/{id}/places",
		Summary:       "List places in a location",
		Tags:          []string{"Places"},
		DefaultStatus: http.StatusOK,
	}, h.HandleListByLocation)

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

func (h *PlaceHandler) HandleGetByID(ctx context.Context, in *GetPlaceRequest) (*PlaceWithRatingResponse, error) {
	p, err := h.svc.GetByIDWithRating(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &PlaceWithRatingResponse{Body: p}, nil
}

func (h *PlaceHandler) HandleListByLocation(ctx context.Context, in *ListPlacesRequest) (*ListPlacesResponse, error) {
	items, err := h.svc.ListByLocation(ctx, in.ID, in.Category, placesvc.SortKey(in.Sort))
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &ListPlacesResponse{Body: items}, nil
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
			Latitude:  in.Body.Coordinates.Latitude,
			Longitude: in.Body.Coordinates.Longitude,
		},
		Description: in.Body.Description,
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
	}
	if in.Body.Coordinates != nil {
		updateInput.Coordinates = &coordinates.Coordinates{
			Latitude:  in.Body.Coordinates.Latitude,
			Longitude: in.Body.Coordinates.Longitude,
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
