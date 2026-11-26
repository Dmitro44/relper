package ratings

import (
	"context"
	"errors"
	"testing"

	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/rating"
	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/google/uuid"
)

type fakeRepo struct {
	upsertedUserID  uuid.UUID
	upsertedPlaceID uuid.UUID
	upsertedStars   int16
	upsertCalls     int
	upsertErr       error

	deleteCalls int
	deleteErr   error
}

func (f *fakeRepo) Upsert(ctx context.Context, userID, placeID uuid.UUID, stars int16) (*rating.PlaceRating, error) {
	f.upsertCalls++
	f.upsertedUserID = userID
	f.upsertedPlaceID = placeID
	f.upsertedStars = stars
	if f.upsertErr != nil {
		return nil, f.upsertErr
	}
	return &rating.PlaceRating{UserID: userID, PlaceID: placeID, Stars: stars}, nil
}

func (f *fakeRepo) Delete(ctx context.Context, userID, placeID uuid.UUID) error {
	f.deleteCalls++
	return f.deleteErr
}

func (f *fakeRepo) Get(ctx context.Context, userID, placeID uuid.UUID) (*rating.PlaceRating, error) {
	return &rating.PlaceRating{UserID: userID, PlaceID: placeID, Stars: 4}, nil
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

func TestSet(t *testing.T) {
	userID, placeID := uuid.New(), uuid.New()

	for _, stars := range []int16{0, 6, -1} {
		t.Run("stars out of range", func(t *testing.T) {
			repo := &fakeRepo{}
			svc := NewRatingService(repo, &fakePlaceGetter{})

			_, err := svc.Set(context.Background(), userID, placeID, stars)

			assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
			if repo.upsertCalls != 0 {
				t.Fatal("expected no upsert")
			}
		})
	}

	t.Run("place not found", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewRatingService(repo, &fakePlaceGetter{err: apperr.NewNotFoundError("place not found")})

		_, err := svc.Set(context.Background(), userID, placeID, 4)

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
		if repo.upsertCalls != 0 {
			t.Fatal("expected no upsert")
		}
	})

	t.Run("sets rating", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewRatingService(repo, &fakePlaceGetter{})

		got, err := svc.Set(context.Background(), userID, placeID, 4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.upsertCalls != 1 {
			t.Fatalf("expected 1 upsert, got %d", repo.upsertCalls)
		}
		if repo.upsertedUserID != userID || repo.upsertedPlaceID != placeID || repo.upsertedStars != 4 {
			t.Fatalf("unexpected upsert args: %s %s %d", repo.upsertedUserID, repo.upsertedPlaceID, repo.upsertedStars)
		}
		if got == nil || got.Stars != 4 {
			t.Fatalf("unexpected rating: %+v", got)
		}
	})

	t.Run("upsert failure", func(t *testing.T) {
		repo := &fakeRepo{upsertErr: errors.New("boom")}
		svc := NewRatingService(repo, &fakePlaceGetter{})

		_, err := svc.Set(context.Background(), userID, placeID, 3)

		assertErrorCode(t, err, apperr.ErrorCodeInternal)
	})
}

func TestRemove(t *testing.T) {
	userID, placeID := uuid.New(), uuid.New()

	t.Run("idempotent", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewRatingService(repo, &fakePlaceGetter{})

		if err := svc.Remove(context.Background(), userID, placeID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := svc.Remove(context.Background(), userID, placeID); err != nil {
			t.Fatalf("unexpected error on second remove: %v", err)
		}
		if repo.deleteCalls != 2 {
			t.Fatalf("expected 2 deletes, got %d", repo.deleteCalls)
		}
	})

	t.Run("delete failure", func(t *testing.T) {
		repo := &fakeRepo{deleteErr: errors.New("boom")}
		svc := NewRatingService(repo, &fakePlaceGetter{})

		err := svc.Remove(context.Background(), userID, placeID)

		assertErrorCode(t, err, apperr.ErrorCodeInternal)
	})
}
