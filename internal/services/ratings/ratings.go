package ratings

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/rating"
	apperr "github.com/Dmitro44/relper/internal/errors"
	ratingsrepo "github.com/Dmitro44/relper/internal/repositories/ratings"
	"github.com/google/uuid"
)

type PlaceGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error)
}

type RatingService struct {
	repo   ratingsrepo.Repository
	places PlaceGetter
}

func NewRatingService(repo ratingsrepo.Repository, places PlaceGetter) *RatingService {
	return &RatingService{
		repo:   repo,
		places: places,
	}
}

func (s *RatingService) Set(ctx context.Context, userID, placeID uuid.UUID, stars int16) (*rating.PlaceRating, error) {
	if stars < 1 || stars > 5 {
		return nil, apperr.NewBadRequestError("stars must be between 1 and 5")
	}

	if _, err := s.places.GetByID(ctx, placeID); err != nil {
		return nil, err
	}

	placeRating, err := s.repo.Upsert(ctx, userID, placeID, stars)
	if err != nil {
		return nil, apperr.NewInternalError("failed to set place rating", err)
	}
	slog.InfoContext(ctx, fmt.Sprintf("place rating set to %d stars", stars),
		"place_id", placeID, "user_id", userID)

	return placeRating, nil
}

func (s *RatingService) Remove(ctx context.Context, userID, placeID uuid.UUID) error {
	if err := s.repo.Delete(ctx, userID, placeID); err != nil {
		return apperr.NewInternalError("failed to remove place rating", err)
	}
	slog.InfoContext(ctx, "place rating removed", "place_id", placeID, "user_id", userID)

	return nil
}

func (s *RatingService) Get(ctx context.Context, userID, placeID uuid.UUID) (*rating.PlaceRating, error) {
	placeRating, err := s.repo.Get(ctx, userID, placeID)
	if err != nil {
		return nil, apperr.NewInternalError("failed to get place rating", err)
	}

	return placeRating, nil
}
