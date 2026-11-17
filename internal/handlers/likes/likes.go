package likes

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	likessvc "github.com/Dmitro44/relper/internal/services/likes"
	"github.com/danielgtaylor/huma/v2"
)

type LikeHandler struct {
	svc *likessvc.LikeService
}

func NewLikeHandler(svc *likessvc.LikeService) *LikeHandler {
	return &LikeHandler{
		svc: svc,
	}
}

func (h *LikeHandler) Register(public, protected huma.API) {
	huma.Register(protected, huma.Operation{
		OperationID:   "like-tip",
		Method:        http.MethodPut,
		Path:          "/tips/{id}/like",
		Summary:       "Like a tip",
		Tags:          []string{"Likes"},
		DefaultStatus: http.StatusOK,
	}, h.HandleLikeTip)

	huma.Register(protected, huma.Operation{
		OperationID:   "unlike-tip",
		Method:        http.MethodDelete,
		Path:          "/tips/{id}/like",
		Summary:       "Unlike a tip",
		Tags:          []string{"Likes"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUnlikeTip)

	huma.Register(public, huma.Operation{
		OperationID:   "get-tip-likes",
		Method:        http.MethodGet,
		Path:          "/tips/{id}/likes",
		Summary:       "Get tip likes",
		Tags:          []string{"Likes"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetTipLikes)

	huma.Register(protected, huma.Operation{
		OperationID:   "like-comment",
		Method:        http.MethodPut,
		Path:          "/comments/{id}/like",
		Summary:       "Like a comment",
		Tags:          []string{"Likes"},
		DefaultStatus: http.StatusOK,
	}, h.HandleLikeComment)

	huma.Register(protected, huma.Operation{
		OperationID:   "unlike-comment",
		Method:        http.MethodDelete,
		Path:          "/comments/{id}/like",
		Summary:       "Unlike a comment",
		Tags:          []string{"Likes"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUnlikeComment)

	huma.Register(public, huma.Operation{
		OperationID:   "get-comment-likes",
		Method:        http.MethodGet,
		Path:          "/comments/{id}/likes",
		Summary:       "Get comment likes",
		Tags:          []string{"Likes"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetCommentLikes)
}

func (h *LikeHandler) HandleLikeTip(ctx context.Context, in *LikeTipRequest) (*SummaryResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	summary, err := h.svc.LikeTip(ctx, userID, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SummaryResponse{Body: summary}, nil
}

func (h *LikeHandler) HandleUnlikeTip(ctx context.Context, in *UnlikeTipRequest) (*SummaryResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	summary, err := h.svc.UnlikeTip(ctx, userID, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SummaryResponse{Body: summary}, nil
}

func (h *LikeHandler) HandleGetTipLikes(ctx context.Context, in *GetTipLikesRequest) (*SummaryResponse, error) {
	summary, err := h.svc.TipSummary(ctx, in.ID, nil)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SummaryResponse{Body: summary}, nil
}

func (h *LikeHandler) HandleLikeComment(ctx context.Context, in *LikeCommentRequest) (*SummaryResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	summary, err := h.svc.LikeComment(ctx, userID, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SummaryResponse{Body: summary}, nil
}

func (h *LikeHandler) HandleUnlikeComment(ctx context.Context, in *UnlikeCommentRequest) (*SummaryResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	summary, err := h.svc.UnlikeComment(ctx, userID, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SummaryResponse{Body: summary}, nil
}

func (h *LikeHandler) HandleGetCommentLikes(ctx context.Context, in *GetCommentLikesRequest) (*SummaryResponse, error) {
	summary, err := h.svc.CommentSummary(ctx, in.ID, nil)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SummaryResponse{Body: summary}, nil
}
