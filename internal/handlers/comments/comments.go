package comments

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/authn"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	commentsvc "github.com/Dmitro44/relper/internal/services/comments"
	"github.com/danielgtaylor/huma/v2"
)

type CommentHandler struct {
	svc *commentsvc.CommentService
}

func NewCommentHandler(svc *commentsvc.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

func (h *CommentHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Comments"}
	})

	huma.Get(grp, "/comments/{id}", h.HandleGetByID)
	huma.Patch(grp, "/comments/{id}", h.HandleUpdate)
	huma.Delete(grp, "/comments/{id}", h.HandleDelete)

	huma.Register(grp, huma.Operation{
		OperationID:   "create-comment",
		Method:        http.MethodPost,
		Path:          "/comments",
		Summary:       "Create a comment",
		DefaultStatus: 201,
	}, h.HandleCreate)
}

func (h *CommentHandler) HandleGetByID(ctx context.Context, in *GetCommentRequest) (*CommentResponse, error) {
	u, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentResponse{Body: u}, nil
}

func (h *CommentHandler) HandleCreate(ctx context.Context, in *CreateCommentRequest) (*CommentResponse, error) {
	userID, ok := authn.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	u, err := h.svc.Create(ctx, commentsvc.CreateInput{
		TargetType: in.Body.TargetType,
		TargetID:   in.Body.TargetID,
		UserID:     userID,
		Content:    in.Body.Content,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentResponse{Body: u}, nil
}

func (h *CommentHandler) HandleUpdate(ctx context.Context, in *UpdateCommentRequest) (*CommentResponse, error) {
	userID, ok := authn.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	updateInput := commentsvc.UpdateInput{
		TargetType: in.Body.TargetType,
		TargetID:   in.Body.TargetID,
		Content:    in.Body.Content,
	}

	u, err := h.svc.Update(ctx, in.ID, userID, updateInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentResponse{Body: u}, nil
}

func (h *CommentHandler) HandleDelete(ctx context.Context, in *DeleteCommentRequest) (*struct{}, error) {
	userID, ok := authn.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	err := h.svc.Delete(ctx, in.ID, userID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
