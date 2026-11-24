package ratings

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	ratingsvc "github.com/Dmitro44/relper/internal/services/ratings"
	"github.com/danielgtaylor/huma/v2"
)

type RatingHandler struct {
	svc *ratingsvc.RatingService
}

func NewRatingHandler(svc *ratingsvc.RatingService) *RatingHandler {
	return &RatingHandler{
		svc: svc,
	}
}

func (h *RatingHandler) Register(public, protected huma.API) {
	huma.Register(protected, huma.Operation{
		OperationID:   "set-place-rating",
		Method:        http.MethodPut,
		Path:          "/places/{id}/rating",
		Summary:       "Set a place rating",
		Tags:          []string{"Ratings"},
		DefaultStatus: http.StatusOK,
	}, h.HandleSet)

	huma.Register(protected, huma.Operation{
		OperationID:   "delete-place-rating",
		Method:        http.MethodDelete,
		Path:          "/places/{id}/rating",
		Summary:       "Remove a place rating",
		Tags:          []string{"Ratings"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)
}

func (h *RatingHandler) HandleSet(ctx context.Context, in *SetRatingRequest) (*RatingResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	placeRating, err := h.svc.Set(ctx, userID, in.ID, in.Body.Stars)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &RatingResponse{Body: toRatingItem(*placeRating)}, nil
}

func (h *RatingHandler) HandleDelete(ctx context.Context, in *DeleteRatingRequest) (*struct{}, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	if err := h.svc.Remove(ctx, userID, in.ID); err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
