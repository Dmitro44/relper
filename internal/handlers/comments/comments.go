package comments

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
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

func (h *CommentHandler) Register(public, protected huma.API) {
	huma.Register(public, huma.Operation{
		OperationID:   "get-comment",
		Method:        http.MethodGet,
		Path:          "/comments/{id}",
		Summary:       "Get a comment",
		Tags:          []string{"Comments"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetByID)

	huma.Register(public, huma.Operation{
		OperationID:   "list-place-comments",
		Method:        http.MethodGet,
		Path:          "/places/{id}/comments",
		Summary:       "List comments for a place",
		Tags:          []string{"Comments"},
		DefaultStatus: http.StatusOK,
	}, h.HandleListByPlace)

	huma.Register(public, huma.Operation{
		OperationID:   "list-tip-comments",
		Method:        http.MethodGet,
		Path:          "/tips/{id}/comments",
		Summary:       "List comments for a tip",
		Tags:          []string{"Comments"},
		DefaultStatus: http.StatusOK,
	}, h.HandleListByTip)

	huma.Register(protected, huma.Operation{
		OperationID:   "create-comment",
		Method:        http.MethodPost,
		Path:          "/comments",
		Summary:       "Create a comment",
		Tags:          []string{"Comments"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)

	huma.Register(protected, huma.Operation{
		OperationID:   "update-comment",
		Method:        http.MethodPatch,
		Path:          "/comments/{id}",
		Summary:       "Update a comment",
		Tags:          []string{"Comments"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUpdate)

	huma.Register(protected, huma.Operation{
		OperationID:   "delete-comment",
		Method:        http.MethodDelete,
		Path:          "/comments/{id}",
		Summary:       "Delete a comment",
		Tags:          []string{"Comments"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)
}

func (h *CommentHandler) HandleGetByID(ctx context.Context, in *GetCommentRequest) (*CommentResponse, error) {
	u, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentResponse{Body: toCommentDetail(*u)}, nil
}

func (h *CommentHandler) HandleListByPlace(ctx context.Context, in *ListCommentsRequest) (*CommentsResponse, error) {
	items, err := h.svc.ListByPlace(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentsResponse{Body: toCommentItems(items)}, nil
}

func (h *CommentHandler) HandleListByTip(ctx context.Context, in *ListCommentsRequest) (*CommentsResponse, error) {
	items, err := h.svc.ListByTip(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentsResponse{Body: toCommentItems(items)}, nil
}

func (h *CommentHandler) HandleCreate(ctx context.Context, in *CreateCommentRequest) (*CommentResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	u, err := h.svc.Create(ctx, commentsvc.CreateInput{
		PlaceID: in.Body.PlaceID,
		TipID:   in.Body.TipID,
		UserID:  userID,
		Content: in.Body.Content,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentResponse{Body: toCommentDetail(*u)}, nil
}

func (h *CommentHandler) HandleUpdate(ctx context.Context, in *UpdateCommentRequest) (*CommentResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	updateInput := commentsvc.UpdateInput{
		Content: in.Body.Content,
	}

	u, err := h.svc.Update(ctx, in.ID, userID, auth.HasRole(ctx, "admin"), updateInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &CommentResponse{Body: toCommentDetail(*u)}, nil
}

func (h *CommentHandler) HandleDelete(ctx context.Context, in *DeleteCommentRequest) (*struct{}, error) {
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
