package locations

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
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

func (h *LocationHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Locations"}
	})

	huma.Get(grp, "/locations/{id}", h.HandleGetByID)
	huma.Patch(grp, "/locations/{id}", h.HandleUpdate)
	huma.Delete(grp, "/locations/{id}", h.HandleDelete)

	huma.Register(grp, huma.Operation{
		OperationID:   "create-location",
		Method:        http.MethodPost,
		Path:          "/locations",
		Summary:       "Create a location",
		DefaultStatus: 201,
	}, h.HandleCreate)
}

func (h *LocationHandler) HandleGetByID(ctx context.Context, in *GetLocationRequest) (*LocationResponse, error) {
	loc, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return &LocationResponse{Body: loc}, nil
}

func (h *LocationHandler) HandleDelete(ctx context.Context, in *DeleteLocationRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
