package tips

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	tipsvc "github.com/Dmitro44/relper/internal/services/tips"
	"github.com/danielgtaylor/huma/v2"
)

type TipHandler struct {
	svc *tipsvc.TipService
}

func NewTipHandler(svc *tipsvc.TipService) *TipHandler {
	return &TipHandler{
		svc: svc,
	}
}

func (h *TipHandler) Register(public, protected huma.API) {
	huma.Register(public, huma.Operation{
		OperationID:   "get-tip",
		Method:        http.MethodGet,
		Path:          "/tips/{id}",
		Summary:       "Get a tip",
		Tags:          []string{"Tips"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetByID)

	huma.Register(public, huma.Operation{
		OperationID:   "list-place-tips",
		Method:        http.MethodGet,
		Path:          "/places/{id}/tips",
		Summary:       "List tips for a place",
		Tags:          []string{"Tips"},
		DefaultStatus: http.StatusOK,
	}, h.HandleListByPlace)

	huma.Register(protected, huma.Operation{
		OperationID:   "create-tip",
		Method:        http.MethodPost,
		Path:          "/tips",
		Summary:       "Create a tip",
		Tags:          []string{"Tips"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)

	huma.Register(protected, huma.Operation{
		OperationID:   "update-tip",
		Method:        http.MethodPatch,
		Path:          "/tips/{id}",
		Summary:       "Update a tip",
		Tags:          []string{"Tips"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUpdate)

	huma.Register(protected, huma.Operation{
		OperationID:   "delete-tip",
		Method:        http.MethodDelete,
		Path:          "/tips/{id}",
		Summary:       "Delete a tip",
		Tags:          []string{"Tips"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)
}

func (h *TipHandler) HandleGetByID(ctx context.Context, in *GetTipRequest) (*TipResponse, error) {
	t, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleListByPlace(ctx context.Context, in *ListPlaceTipsRequest) (*TipsResponse, error) {
	items, err := h.svc.ListByPlace(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TipsResponse{Body: items}, nil
}

func (h *TipHandler) HandleCreate(ctx context.Context, in *CreateTipRequest) (*TipResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	t, err := h.svc.Create(ctx, tipsvc.CreateInput{
		LocationID: in.Body.LocationID,
		PlaceID:    in.Body.PlaceID,
		UserID:     userID,
		Category:   in.Body.Category,
		Title:      in.Body.Title,
		Content:    in.Body.Content,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleUpdate(ctx context.Context, in *UpdateTipRequest) (*TipResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	t, err := h.svc.Update(ctx, in.ID, userID, auth.HasRole(ctx, "admin"), tipsvc.UpdateInput{
		LocationID: in.Body.LocationID,
		PlaceID:    in.Body.PlaceID,
		Category:   in.Body.Category,
		Title:      in.Body.Title,
		Content:    in.Body.Content,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleDelete(ctx context.Context, in *DeleteTipRequest) (*struct{}, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	err := h.svc.Delete(ctx, in.ID, userID, auth.HasRole(ctx, "admin"))
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
