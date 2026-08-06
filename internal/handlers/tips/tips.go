package tips

import (
	"context"
	"net/http"

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

func (h *TipHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Tips"}
	})

	huma.Get(grp, "/tips/{id}", h.HandleGetByID)
	huma.Patch(grp, "/tips/{id}", h.HandleUpdate)
	huma.Delete(grp, "/tips/{id}", h.HandleDelete)
	huma.Post(grp, "/tips/{id}/upvote", h.HandleUpvote, func(o *huma.Operation) {
		o.Summary = "Add an upvote"
	})

	huma.Register(grp, huma.Operation{
		OperationID:   "remove-tip-upvote",
		Method:        http.MethodDelete,
		Path:          "/tips/{id}/upvote",
		Summary:       "Remove upvote",
		DefaultStatus: 200,
	}, h.HandleRemoveUpvote)

	huma.Register(grp, huma.Operation{
		OperationID:   "create-tip",
		Method:        http.MethodPost,
		Path:          "/tips",
		Summary:       "Create a tip",
		DefaultStatus: 201,
	}, h.HandleCreate)
}

func (h *TipHandler) HandleGetByID(ctx context.Context, in *GetTipRequest) (*TipResponse, error) {
	t, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleCreate(ctx context.Context, in *CreateTipRequest) (*TipResponse, error) {
	t, err := h.svc.Create(ctx, tipsvc.CreateInput{
		LocationID: in.Body.LocationID,
		UserID:     in.Body.UserID,
		Category:   in.Body.Category,
		Title:      in.Body.Title,
		Content:    in.Body.Content,
		Upvotes:    in.Body.Upvotes,
	})
	if err != nil {
		return nil, err
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleUpdate(ctx context.Context, in *UpdateTipRequest) (*TipResponse, error) {
	t, err := h.svc.Update(ctx, in.ID, tipsvc.UpdateInput{
		LocationID: in.Body.LocationID,
		UserID:     in.Body.UserID,
		Category:   in.Body.Category,
		Title:      in.Body.Title,
		Content:    in.Body.Content,
		Upvotes:    in.Body.Upvotes,
	})
	if err != nil {
		return nil, err
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleDelete(ctx context.Context, in *DeleteTipRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *TipHandler) HandleUpvote(ctx context.Context, in *UpvoteTipRequest) (*TipResponse, error) {
	t, err := h.svc.Upvote(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &TipResponse{Body: t}, nil
}

func (h *TipHandler) HandleRemoveUpvote(ctx context.Context, in *RemoveUpvoteTipRequest) (*TipResponse, error) {
	t, err := h.svc.RemoveUpvote(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &TipResponse{Body: t}, nil
}
