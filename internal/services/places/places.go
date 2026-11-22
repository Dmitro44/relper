package places

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/rating"
	apperr "github.com/Dmitro44/relper/internal/errors"
	placesrepo "github.com/Dmitro44/relper/internal/repositories/places"
	"github.com/google/uuid"
)

type LocationGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type SortKey string

const (
	SortWilson   SortKey = "wilson"
	SortAvgStars SortKey = "avg_stars"
	SortVotes    SortKey = "votes_count"
	SortName     SortKey = "name"
)

// primary returns 0 when the sort key does not decide the order.
func primary(a, b place.WithRating, key SortKey) int {
	switch key {
	case SortAvgStars:
		return cmp.Compare(b.Rating.AvgStars, a.Rating.AvgStars)
	case SortVotes:
		return cmp.Compare(b.Rating.VotesCount, a.Rating.VotesCount)
	case SortName:
		return strings.Compare(a.Place.Name, b.Place.Name)
	default:
		return cmp.Compare(b.Rating.Wilson, a.Rating.Wilson)
	}
}

type PlaceService struct {
	repo            placesrepo.Repository
	locationservice LocationGetter
}

func NewPlaceService(repo placesrepo.Repository, locationservice LocationGetter) *PlaceService {
	return &PlaceService{
		repo:            repo,
		locationservice: locationservice,
	}
}

type CreateInput struct {
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Coordinates coordinates.Coordinates
	Description string
}

type UpdateInput struct {
	LocationID  *uuid.UUID
	Name        *string
	Category    *string
	Address     *string
	Phone       *string
	WebsiteURL  *string
	Coordinates *coordinates.Coordinates
	Description *string
}

func (p *PlaceService) GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error) {
	placeEntity, err := p.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, placesrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("place not found")
		}
		return nil, apperr.NewInternalError("failed to get place", err)
	}

	return placeEntity, nil
}

func (p *PlaceService) ListByLocation(ctx context.Context, locationID uuid.UUID, category string, sortBy SortKey) ([]place.WithRating, error) {
	if _, err := p.locationservice.GetByID(ctx, locationID); err != nil {
		return nil, err
	}

	items, err := p.repo.ListByLocation(ctx, locationID, category)
	if err != nil {
		return nil, apperr.NewInternalError("failed to list places", err)
	}

	result := make([]place.WithRating, 0, len(items))
	for _, item := range items {
		result = append(result, toPlaceWithRating(item))
	}

	slices.SortFunc(result, func(a, b place.WithRating) int {
		return cmp.Or(
			primary(a, b, sortBy),
			cmp.Compare(b.Rating.VotesCount, a.Rating.VotesCount),
			b.Place.CreatedAt.Compare(a.Place.CreatedAt),
			strings.Compare(a.Place.ID.String(), b.Place.ID.String()),
		)
	})

	return result, nil
}

func (p *PlaceService) GetByIDWithRating(ctx context.Context, id uuid.UUID) (*place.WithRating, error) {
	item, err := p.repo.GetByIDWithRating(ctx, id)
	if err != nil {
		if errors.Is(err, placesrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("place not found")
		}
		return nil, apperr.NewInternalError("failed to get place", err)
	}

	result := toPlaceWithRating(*item)

	return &result, nil
}

func toPlaceWithRating(item placesrepo.PlaceWithRating) place.WithRating {
	return place.WithRating{
		Place: item.Place,
		Rating: rating.Summary{
			AvgStars:   item.RatingCounts.AvgStars,
			VotesCount: item.RatingCounts.VotesCount,
			Wilson:     rating.Wilson(item.RatingCounts.Recommend, item.RatingCounts.NotRecommend, 1.96),
		},
	}
}

func (p *PlaceService) Create(ctx context.Context, input CreateInput) (*place.Place, error) {
	// check if location id exists
	_, err := p.locationservice.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, err
	}

	placeEntity, err := p.repo.Create(ctx, placesrepo.CreateInput{
		LocationID:  input.LocationID,
		Name:        input.Name,
		Category:    input.Category,
		Address:     input.Address,
		Phone:       input.Phone,
		WebsiteURL:  input.WebsiteURL,
		Coordinates: input.Coordinates,
		Description: input.Description,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create place", err)
	}
	slog.InfoContext(ctx, "place created", "place_id", placeEntity.ID)

	return placeEntity, nil
}

func (p *PlaceService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*place.Place, error) {
	placeEntity, err := p.repo.Update(ctx, id, placesrepo.UpdateInput{
		LocationID:  input.LocationID,
		Name:        input.Name,
		Category:    input.Category,
		Address:     input.Address,
		Phone:       input.Phone,
		WebsiteURL:  input.WebsiteURL,
		Coordinates: input.Coordinates,
		Description: input.Description,
	})
	if err != nil {
		if errors.Is(err, placesrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("place not found")
		}
		return nil, apperr.NewInternalError("failed to update place", err)
	}
	slog.InfoContext(ctx, "place updated", "place_id", placeEntity.ID)

	return placeEntity, nil
}

func (p *PlaceService) Delete(ctx context.Context, id uuid.UUID) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, placesrepo.ErrNotFound) {
			return apperr.NewNotFoundError("place not found")
		}
		return apperr.NewInternalError("failed to delete place", err)
	}
	slog.InfoContext(ctx, "place deleted", "place_id", id)

	return nil
}
