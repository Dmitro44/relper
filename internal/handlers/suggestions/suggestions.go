package suggestions

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	suggestionsvc "github.com/Dmitro44/relper/internal/services/suggestions"
	"github.com/danielgtaylor/huma/v2"
)

type SuggestionHandler struct {
	svc *suggestionsvc.SuggestionService
}

func NewSuggestionHandler(svc *suggestionsvc.SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{svc: svc}
}

func (h *SuggestionHandler) Register(protected, admin huma.API) {
	huma.Register(protected, huma.Operation{
		OperationID:   "create-place-suggestion",
		Method:        http.MethodPost,
		Path:          "/place-suggestions",
		Summary:       "Suggest a new place",
		Tags:          []string{"Suggestions"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)

	huma.Register(protected, huma.Operation{
		OperationID:   "list-my-place-suggestions",
		Method:        http.MethodGet,
		Path:          "/place-suggestions/mine",
		Summary:       "List the current user's place suggestions",
		Tags:          []string{"Suggestions"},
		DefaultStatus: http.StatusOK,
	}, h.HandleListMine)

	huma.Register(admin, huma.Operation{
		OperationID:   "list-place-suggestions",
		Method:        http.MethodGet,
		Path:          "/place-suggestions",
		Summary:       "List place suggestions for review",
		Tags:          []string{"Suggestions"},
		DefaultStatus: http.StatusOK,
	}, h.HandleList)

	huma.Register(admin, huma.Operation{
		OperationID:   "approve-place-suggestion",
		Method:        http.MethodPatch,
		Path:          "/place-suggestions/{id}/approve",
		Summary:       "Approve a place suggestion",
		Tags:          []string{"Suggestions"},
		DefaultStatus: http.StatusOK,
	}, h.HandleApprove)

	huma.Register(admin, huma.Operation{
		OperationID:   "reject-place-suggestion",
		Method:        http.MethodPatch,
		Path:          "/place-suggestions/{id}/reject",
		Summary:       "Reject a place suggestion",
		Tags:          []string{"Suggestions"},
		DefaultStatus: http.StatusOK,
	}, h.HandleReject)
}

func (h *SuggestionHandler) HandleCreate(ctx context.Context, in *CreateSuggestionRequest) (*SuggestionResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	item, err := h.svc.Create(ctx, suggestionsvc.CreateInput{
		UserID:      userID,
		LocationID:  in.Body.LocationID,
		Name:        in.Body.Name,
		Category:    in.Body.Category,
		Address:     in.Body.Address,
		Phone:       in.Body.Phone,
		WebsiteURL:  in.Body.WebsiteURL,
		Description: in.Body.Description,
		Coordinates: coordinates.Coordinates{
			Latitude:  in.Body.Coordinates.Latitude,
			Longitude: in.Body.Coordinates.Longitude,
		},
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SuggestionResponse{Body: toSuggestionItem(*item)}, nil
}

func (h *SuggestionHandler) HandleListMine(ctx context.Context, _ *struct{}) (*ListSuggestionsResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	items, err := h.svc.ListMine(ctx, userID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &ListSuggestionsResponse{Body: toSuggestionItems(items)}, nil
}

func (h *SuggestionHandler) HandleList(ctx context.Context, in *ListSuggestionsRequest) (*ListSuggestionsResponse, error) {
	items, err := h.svc.List(ctx, in.Status)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &ListSuggestionsResponse{Body: toSuggestionItems(items)}, nil
}

func (h *SuggestionHandler) HandleApprove(ctx context.Context, in *ApproveSuggestionRequest) (*SuggestionResponse, error) {
	reviewerID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	item, err := h.svc.Approve(ctx, in.ID, reviewerID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SuggestionResponse{Body: toSuggestionItem(*item)}, nil
}

func (h *SuggestionHandler) HandleReject(ctx context.Context, in *RejectSuggestionRequest) (*SuggestionResponse, error) {
	reviewerID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	item, err := h.svc.Reject(ctx, in.ID, reviewerID, in.Body.Reason)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &SuggestionResponse{Body: toSuggestionItem(*item)}, nil
}
