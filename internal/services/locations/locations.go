package locations

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/location"
	apperr "github.com/Dmitro44/relper/internal/errors"
	locationsrepo "github.com/Dmitro44/relper/internal/repositories/locations"
	"github.com/google/uuid"
)

type LocationService struct {
	repo locationsrepo.Repository
}

type CreateInput struct {
	Country     string
	City        string
	Coordinates coordinates.Coordinates
}
type UpdateInput struct {
	Country     *string
	City        *string
	Coordinates *coordinates.Coordinates
}

func NewLocationService(repo locationsrepo.Repository) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (l *LocationService) GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error) {
	locationEntity, err := l.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, locationsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("location not found")
		}
		return nil, apperr.NewInternalError("failed to get location", err)
	}

	return locationEntity, nil
}

func (l *LocationService) Create(ctx context.Context, input CreateInput) (*location.Location, error) {
	locationEntity, err := l.repo.Create(ctx, locationsrepo.CreateInput{
		Country:     input.Country,
		City:        input.City,
		Coordinates: input.Coordinates,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create location", err)
	}
	slog.InfoContext(ctx, "location created", "location_id", locationEntity.ID)

	return locationEntity, nil
}

func (l *LocationService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*location.Location, error) {
	locationEntity, err := l.repo.Update(ctx, id, locationsrepo.UpdateInput{
		Country:     input.Country,
		City:        input.City,
		Coordinates: input.Coordinates,
	})
	if err != nil {
		if errors.Is(err, locationsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("location not found")
		}
		return nil, apperr.NewInternalError("failed to update location", err)
	}
	slog.InfoContext(ctx, "location updated", "location_id", locationEntity.ID)

	return locationEntity, nil
}

func (l *LocationService) Delete(ctx context.Context, id uuid.UUID) error {
	err := l.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, locationsrepo.ErrNotFound) {
			return apperr.NewNotFoundError("location not found")
		}
		return apperr.NewInternalError("failed to delete location", err)
	}
	slog.InfoContext(ctx, "location deleted", "location_id", id)

	return nil
}
