package suggestions

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/suggestion"
	apperr "github.com/Dmitro44/relper/internal/errors"
	suggestionsrepo "github.com/Dmitro44/relper/internal/repositories/suggestions"
	"github.com/google/uuid"
)

type LocationGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type SuggestionService struct {
	repo      suggestionsrepo.Repository
	locations LocationGetter
}

func NewSuggestionService(repo suggestionsrepo.Repository, locations LocationGetter) *SuggestionService {
	return &SuggestionService{repo: repo, locations: locations}
}

type CreateInput struct {
	UserID      uuid.UUID
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Description string
	Coordinates coordinates.Coordinates
}

func (s *SuggestionService) Create(ctx context.Context, input CreateInput) (*suggestion.Suggestion, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Category = strings.TrimSpace(input.Category)
	if input.Name == "" {
		return nil, apperr.NewBadRequestError("name is required")
	}
	if !place.IsValidCategory(input.Category) {
		return nil, apperr.NewBadRequestError("unknown place category")
	}
	if input.Coordinates.Latitude < -90 || input.Coordinates.Latitude > 90 ||
		input.Coordinates.Longitude < -180 || input.Coordinates.Longitude > 180 {
		return nil, apperr.NewBadRequestError("invalid coordinates")
	}

	if _, err := s.locations.GetByID(ctx, input.LocationID); err != nil {
		return nil, err
	}

	item, err := s.repo.Create(ctx, suggestionsrepo.CreateInput(input))
	if err != nil {
		return nil, apperr.NewInternalError("failed to create suggestion", err)
	}
	slog.InfoContext(ctx, "place suggestion created", "suggestion_id", item.ID, "user_id", input.UserID)

	return item, nil
}

func (s *SuggestionService) ListMine(ctx context.Context, userID uuid.UUID) ([]suggestion.Suggestion, error) {
	items, err := s.repo.ListMine(ctx, userID)
	if err != nil {
		return nil, apperr.NewInternalError("failed to list suggestions", err)
	}

	return items, nil
}

func (s *SuggestionService) List(ctx context.Context, status string) ([]suggestion.Suggestion, error) {
	if status == "" {
		status = suggestion.StatusPending
	}
	switch status {
	case suggestion.StatusPending, suggestion.StatusApproved, suggestion.StatusRejected:
	default:
		return nil, apperr.NewBadRequestError("unknown suggestion status")
	}

	items, err := s.repo.ListByStatus(ctx, status)
	if err != nil {
		return nil, apperr.NewInternalError("failed to list suggestions", err)
	}

	return items, nil
}

func (s *SuggestionService) Approve(ctx context.Context, id, reviewerID uuid.UUID) (*suggestion.Suggestion, error) {
	item, err := s.repo.Approve(ctx, id, reviewerID)
	if err != nil {
		switch {
		case errors.Is(err, suggestionsrepo.ErrNotFound):
			return nil, apperr.NewNotFoundError("suggestion not found")
		case errors.Is(err, suggestionsrepo.ErrNotPending):
			return nil, apperr.NewConflictError("suggestion is already reviewed")
		}
		return nil, apperr.NewInternalError("failed to approve suggestion", err)
	}
	slog.InfoContext(ctx, "place suggestion approved", "suggestion_id", id, "reviewer_id", reviewerID)

	return item, nil
}

func (s *SuggestionService) Reject(ctx context.Context, id, reviewerID uuid.UUID, reason string) (*suggestion.Suggestion, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, apperr.NewBadRequestError("rejection reason is required")
	}

	item, err := s.repo.Reject(ctx, id, reviewerID, reason)
	if err != nil {
		switch {
		case errors.Is(err, suggestionsrepo.ErrNotFound):
			return nil, apperr.NewNotFoundError("suggestion not found")
		case errors.Is(err, suggestionsrepo.ErrNotPending):
			return nil, apperr.NewConflictError("suggestion is already reviewed")
		}
		return nil, apperr.NewInternalError("failed to reject suggestion", err)
	}
	slog.InfoContext(ctx, "place suggestion rejected", "suggestion_id", id, "reviewer_id", reviewerID)

	return item, nil
}
