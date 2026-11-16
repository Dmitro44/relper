package tips

import (
	"context"
	"errors"
	"testing"

	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips"
	"github.com/google/uuid"
)

type fakeRepo struct {
	created tipsrepo.CreateInput
	updated tipsrepo.UpdateInput
}

func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	return &tip.Tip{ID: id}, nil
}

func (f *fakeRepo) Create(ctx context.Context, input tipsrepo.CreateInput) (*tip.Tip, error) {
	f.created = input
	return &tip.Tip{
		ID:         uuid.New(),
		LocationID: input.LocationID,
		PlaceID:    input.PlaceID,
		UserID:     input.UserID,
		Category:   input.Category,
		Title:      input.Title,
		Content:    input.Content,
	}, nil
}

func (f *fakeRepo) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input tipsrepo.UpdateInput) (*tip.Tip, error) {
	f.updated = input
	return &tip.Tip{ID: id, UserID: userID}, nil
}

func (f *fakeRepo) Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error {
	return nil
}

type fakePlaceGetter struct {
	err error
}

func (f *fakePlaceGetter) GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &place.Place{ID: id}, nil
}

type fakeLocationGetter struct {
	err error
}

func (f *fakeLocationGetter) GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &location.Location{ID: id}, nil
}

func assertErrorCode(t *testing.T, err error, want apperr.ErrorCode) {
	t.Helper()

	var appErr *apperr.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *apperr.AppError, got %v", err)
	}
	if appErr.Code != want {
		t.Fatalf("expected code %s, got %s", want, appErr.Code)
	}
}

func TestCreate(t *testing.T) {
	userID := uuid.New()

	t.Run("both targets", func(t *testing.T) {
		locationID, placeID := uuid.New(), uuid.New()
		svc := NewTipService(&fakeRepo{}, &fakePlaceGetter{}, &fakeLocationGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			LocationID: &locationID,
			PlaceID:    &placeID,
			UserID:     userID,
			Category:   "food",
			Title:      "Great",
			Content:    "Try it",
		})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("neither target", func(t *testing.T) {
		svc := NewTipService(&fakeRepo{}, &fakePlaceGetter{}, &fakeLocationGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			UserID:   userID,
			Category: "food",
			Title:    "Great",
			Content:  "Try it",
		})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("place not found", func(t *testing.T) {
		placeID := uuid.New()
		svc := NewTipService(
			&fakeRepo{},
			&fakePlaceGetter{err: apperr.NewNotFoundError("place not found")},
			&fakeLocationGetter{},
		)

		_, err := svc.Create(context.Background(), CreateInput{
			PlaceID:  &placeID,
			UserID:   userID,
			Category: "food",
			Title:    "Great",
			Content:  "Try it",
		})

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})

	t.Run("location not found", func(t *testing.T) {
		locationID := uuid.New()
		svc := NewTipService(
			&fakeRepo{},
			&fakePlaceGetter{},
			&fakeLocationGetter{err: apperr.NewNotFoundError("location not found")},
		)

		_, err := svc.Create(context.Background(), CreateInput{
			LocationID: &locationID,
			UserID:     userID,
			Category:   "food",
			Title:      "Great",
			Content:    "Try it",
		})

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("both targets", func(t *testing.T) {
		locationID, placeID := uuid.New(), uuid.New()
		svc := NewTipService(&fakeRepo{}, &fakePlaceGetter{}, &fakeLocationGetter{})

		_, err := svc.Update(context.Background(), uuid.New(), uuid.New(), false, UpdateInput{
			LocationID: &locationID,
			PlaceID:    &placeID,
		})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("place not found", func(t *testing.T) {
		placeID := uuid.New()
		svc := NewTipService(
			&fakeRepo{},
			&fakePlaceGetter{err: apperr.NewNotFoundError("place not found")},
			&fakeLocationGetter{},
		)

		_, err := svc.Update(context.Background(), uuid.New(), uuid.New(), false, UpdateInput{
			PlaceID: &placeID,
		})

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})
}
