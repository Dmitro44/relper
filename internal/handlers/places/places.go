package places

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
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

func (h *PlaceHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Places"}
	})

	huma.Get(grp, "/places/{id}", h.HandleGetByID)
	huma.Patch(grp, "/places/{id}", h.HandleUpdate)
	huma.Delete(grp, "/places/{id}", h.HandleDelete)

	huma.Register(grp, huma.Operation{
		OperationID:   "create-place",
		Method:        http.MethodPost,
		Path:          "/places",
		Summary:       "Create a place",
		DefaultStatus: 201,
	}, h.HandleCreate)
}

func (h *PlaceHandler) HandleGetByID(ctx context.Context, in *GetPlaceRequest) (*PlaceResponse, error) {
	p, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
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
		Upvotes:     in.Body.Upvotes,
	}
	p, err := h.svc.Create(ctx, createInput)
	if err != nil {
		return nil, err
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
		Upvotes:     in.Body.Upvotes,
	}
	if in.Body.Coordinates != nil {
		updateInput.Coordinates = &coordinates.Coordinates{
			Lat: in.Body.Coordinates.Lat,
			Lng: in.Body.Coordinates.Lng,
		}
	}

	p, err := h.svc.Update(ctx, in.ID, updateInput)
	if err != nil {
		return nil, err
	}

	return &PlaceResponse{Body: p}, nil
}

func (h *PlaceHandler) HandleDelete(ctx context.Context, in *DeletePlaceRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
