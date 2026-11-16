package tips

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips"
	"github.com/google/uuid"
)

type PlaceGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error)
}

type LocationGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type TipService struct {
	repo      tipsrepo.Repository
	places    PlaceGetter
	locations LocationGetter
}

type CreateInput struct {
	LocationID *uuid.UUID
	PlaceID    *uuid.UUID
	UserID     uuid.UUID
	Category   string
	Title      string
	Content    string
}

type UpdateInput struct {
	LocationID *uuid.UUID
	PlaceID    *uuid.UUID
	Category   *string
	Title      *string
	Content    *string
}

func NewTipService(repo tipsrepo.Repository, places PlaceGetter, locations LocationGetter) *TipService {
	return &TipService{
		repo:      repo,
		places:    places,
		locations: locations,
	}
}

func (t *TipService) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	tipEntity, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tip not found")
		}
		return nil, apperr.NewInternalError("failed to get tip", err)
	}

	return tipEntity, nil
}

func (t *TipService) Create(ctx context.Context, input CreateInput) (*tip.Tip, error) {
	if (input.LocationID == nil) == (input.PlaceID == nil) {
		return nil, apperr.NewBadRequestError("exactly one of location_id or place_id is required")
	}

	if input.LocationID != nil {
		if _, err := t.locations.GetByID(ctx, *input.LocationID); err != nil {
			return nil, err
		}
	}
	if input.PlaceID != nil {
		if _, err := t.places.GetByID(ctx, *input.PlaceID); err != nil {
			return nil, err
		}
	}

	tipEntity, err := t.repo.Create(ctx, tipsrepo.CreateInput{
		LocationID: input.LocationID,
		PlaceID:    input.PlaceID,
		UserID:     input.UserID,
		Category:   input.Category,
		Title:      input.Title,
		Content:    input.Content,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create tip", err)
	}
	slog.InfoContext(ctx, "tip created", "tip_id", tipEntity.ID)

	return tipEntity, nil
}

func (t *TipService) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input UpdateInput) (*tip.Tip, error) {
	if input.LocationID != nil && input.PlaceID != nil {
		return nil, apperr.NewBadRequestError("only one of location_id or place_id can be set")
	}

	if input.LocationID != nil {
		if _, err := t.locations.GetByID(ctx, *input.LocationID); err != nil {
			return nil, err
		}
	}
	if input.PlaceID != nil {
		if _, err := t.places.GetByID(ctx, *input.PlaceID); err != nil {
			return nil, err
		}
	}

	tipEntity, err := t.repo.Update(ctx, id, userID, admin, tipsrepo.UpdateInput{
		LocationID: input.LocationID,
		PlaceID:    input.PlaceID,
		Category:   input.Category,
		Title:      input.Title,
		Content:    input.Content,
	})
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tip not found")
		}
		return nil, apperr.NewInternalError("failed to update tip", err)
	}
	slog.InfoContext(ctx, "tip updated", "tip_id", tipEntity.ID)

	return tipEntity, nil
}

func (t *TipService) Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error {
	err := t.repo.Delete(ctx, id, userID, admin)
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return apperr.NewNotFoundError("tip not found")
		}
		return apperr.NewInternalError("failed to delete tip", err)
	}
	slog.InfoContext(ctx, "tip deleted", "tip_id", id)

	return nil
}
