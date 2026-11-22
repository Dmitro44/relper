package places

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/rating"
	apperr "github.com/Dmitro44/relper/internal/errors"
	placesrepo "github.com/Dmitro44/relper/internal/repositories/places"
	"github.com/google/uuid"
)

type fakeRepository struct {
	getByIDWithRating    *placesrepo.PlaceWithRating
	getByIDWithRatingErr error
	list                 []placesrepo.PlaceWithRating
	listErr              error
	listCalled           bool
	listLocationID       uuid.UUID
	listCategory         string
}

func (f *fakeRepository) GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error) {
	return &place.Place{ID: id}, nil
}

func (f *fakeRepository) GetByIDWithRating(ctx context.Context, id uuid.UUID) (*placesrepo.PlaceWithRating, error) {
	if f.getByIDWithRatingErr != nil {
		return nil, f.getByIDWithRatingErr
	}
	return f.getByIDWithRating, nil
}

func (f *fakeRepository) ListByLocation(ctx context.Context, locationID uuid.UUID, category string) ([]placesrepo.PlaceWithRating, error) {
	f.listCalled = true
	f.listLocationID = locationID
	f.listCategory = category
	return f.list, f.listErr
}

func (f *fakeRepository) Create(ctx context.Context, input placesrepo.CreateInput) (*place.Place, error) {
	return &place.Place{}, nil
}

func (f *fakeRepository) Update(ctx context.Context, id uuid.UUID, input placesrepo.UpdateInput) (*place.Place, error) {
	return &place.Place{ID: id}, nil
}

func (f *fakeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type fakeLocationGetter struct {
	err    error
	called bool
}

func (f *fakeLocationGetter) GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error) {
	f.called = true
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

func TestListByLocation(t *testing.T) {
	locationID := uuid.New()

	t.Run("location not found", func(t *testing.T) {
		repo := &fakeRepository{}
		locations := &fakeLocationGetter{err: apperr.NewNotFoundError("location not found")}
		svc := NewPlaceService(repo, locations)

		_, err := svc.ListByLocation(context.Background(), locationID, "", SortWilson)

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
		if repo.listCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &fakeRepository{listErr: errors.New("boom")}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		_, err := svc.ListByLocation(context.Background(), locationID, "", SortWilson)

		assertErrorCode(t, err, apperr.ErrorCodeInternal)
	})

	t.Run("default sort uses wilson", func(t *testing.T) {
		repo := &fakeRepository{list: []placesrepo.PlaceWithRating{
			{
				Place:        place.Place{ID: uuid.New(), Name: "Perfect", CreatedAt: time.Now()},
				RatingCounts: placesrepo.RatingCounts{AvgStars: 5.0, VotesCount: 10, Recommend: 10, NotRecommend: 0},
			},
			{
				Place:        place.Place{ID: uuid.New(), Name: "Popular", CreatedAt: time.Now()},
				RatingCounts: placesrepo.RatingCounts{AvgStars: 4.0, VotesCount: 1000, Recommend: 900, NotRecommend: 100},
			},
		}}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		items, err := svc.ListByLocation(context.Background(), locationID, "", SortWilson)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
		if items[0].Place.Name != "Popular" {
			t.Fatalf("expected Popular first, got %s", items[0].Place.Name)
		}

		want := rating.Wilson(900, 100, 1.96)
		if diff := items[0].Rating.Wilson - want; diff > 1e-9 || diff < -1e-9 {
			t.Fatalf("expected wilson %v, got %v", want, items[0].Rating.Wilson)
		}
		if items[0].Rating.VotesCount != 1000 {
			t.Fatalf("expected votes_count 1000, got %d", items[0].Rating.VotesCount)
		}
	})

	t.Run("sort by avg_stars", func(t *testing.T) {
		repo := &fakeRepository{list: []placesrepo.PlaceWithRating{
			{
				Place:        place.Place{ID: uuid.New(), Name: "Popular"},
				RatingCounts: placesrepo.RatingCounts{AvgStars: 4.0, VotesCount: 1000, Recommend: 900, NotRecommend: 100},
			},
			{
				Place:        place.Place{ID: uuid.New(), Name: "Perfect"},
				RatingCounts: placesrepo.RatingCounts{AvgStars: 5.0, VotesCount: 10, Recommend: 10, NotRecommend: 0},
			},
		}}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		items, err := svc.ListByLocation(context.Background(), locationID, "", SortAvgStars)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items[0].Place.Name != "Perfect" {
			t.Fatalf("expected Perfect first, got %s", items[0].Place.Name)
		}
	})

	t.Run("sort by name", func(t *testing.T) {
		repo := &fakeRepository{list: []placesrepo.PlaceWithRating{
			{Place: place.Place{ID: uuid.New(), Name: "Zeta"}},
			{Place: place.Place{ID: uuid.New(), Name: "Alpha"}},
		}}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		items, err := svc.ListByLocation(context.Background(), locationID, "", SortName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items[0].Place.Name != "Alpha" || items[1].Place.Name != "Zeta" {
			t.Fatalf("expected Alpha, Zeta; got %s, %s", items[0].Place.Name, items[1].Place.Name)
		}
	})

	t.Run("category forwarded", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		_, err := svc.ListByLocation(context.Background(), locationID, "cafe", SortWilson)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.listLocationID != locationID {
			t.Fatalf("expected location %s, got %s", locationID, repo.listLocationID)
		}
		if repo.listCategory != "cafe" {
			t.Fatalf("expected category cafe, got %q", repo.listCategory)
		}
	})

	t.Run("empty result is non-nil", func(t *testing.T) {
		svc := NewPlaceService(&fakeRepository{}, &fakeLocationGetter{})

		items, err := svc.ListByLocation(context.Background(), locationID, "", SortWilson)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items == nil {
			t.Fatal("expected non-nil slice")
		}
		if len(items) != 0 {
			t.Fatalf("expected empty slice, got %d", len(items))
		}
	})
}

func TestListByLocationSort(t *testing.T) {
	locationID := uuid.New()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	lowID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	highID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	mk := func(id uuid.UUID, name string, avg float64, votes, recommend, notRecommend int64, createdAt time.Time) placesrepo.PlaceWithRating {
		return placesrepo.PlaceWithRating{
			ID:           id,
			Name:         name,
			CreatedAt:    createdAt,
			AvgStars:     avg,
			VotesCount:   votes,
			Recommend:    recommend,
			NotRecommend: notRecommend,
		}
	}

	ranked := []placesrepo.PlaceWithRating{
		mk(lowID, "Able", 2, 5, 0, 1, base),
		mk(uuid.New(), "Delta", 5, 10, 0, 100, base),
		mk(uuid.New(), "Mike", 3, 50, 50, 0, base),
		mk(highID, "Zulu", 4, 100, 1, 1, base),
	}

	tests := []struct {
		name string
		key  SortKey
		list []placesrepo.PlaceWithRating
		want string
	}{
		{"wilson", SortWilson, ranked, "Mike"},
		{"avg_stars", SortAvgStars, ranked, "Delta"},
		{"votes_count", SortVotes, ranked, "Zulu"},
		{"name", SortName, ranked, "Able"},
		{"empty key falls back to wilson", "", ranked, "Mike"},
		{"unknown key falls back to wilson", SortKey("stars"), ranked, "Mike"},
		{"wilson tie falls through to votes", SortWilson, []placesrepo.PlaceWithRating{
			mk(lowID, "FewVotes", 4, 10, 5, 0, base),
			mk(highID, "ManyVotes", 4, 20, 5, 0, base),
		}, "ManyVotes"},
		{"votes tie falls through to created_at desc", SortWilson, []placesrepo.PlaceWithRating{
			mk(lowID, "Old", 4, 10, 5, 0, base),
			mk(highID, "New", 4, 10, 5, 0, base.Add(time.Hour)),
		}, "New"},
		{"created_at tie falls through to id asc", SortWilson, []placesrepo.PlaceWithRating{
			mk(highID, "HighID", 4, 10, 5, 0, base),
			mk(lowID, "LowID", 4, 10, 5, 0, base),
		}, "LowID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPlaceService(&fakeRepository{list: tt.list}, &fakeLocationGetter{})

			items, err := svc.ListByLocation(context.Background(), locationID, "", tt.key)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(items) == 0 {
				t.Fatal("expected at least one item")
			}
			if items[0].Place.Name != tt.want {
				t.Fatalf("expected %s first, got %s", tt.want, items[0].Place.Name)
			}
		})
	}
}

func TestGetByIDWithRating(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		repo := &fakeRepository{getByIDWithRatingErr: placesrepo.ErrNotFound}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		_, err := svc.GetByIDWithRating(context.Background(), uuid.New())

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})

	t.Run("wilson from recommend counts", func(t *testing.T) {
		repo := &fakeRepository{getByIDWithRating: &placesrepo.PlaceWithRating{
			Place:        place.Place{ID: uuid.New(), Name: "Popular"},
			RatingCounts: placesrepo.RatingCounts{AvgStars: 4.5, VotesCount: 100, Recommend: 10, NotRecommend: 0},
		}}
		svc := NewPlaceService(repo, &fakeLocationGetter{})

		item, err := svc.GetByIDWithRating(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.Rating.AvgStars != 4.5 || item.Rating.VotesCount != 100 {
			t.Fatalf("unexpected rating summary: %+v", item.Rating)
		}

		want := rating.Wilson(10, 0, 1.96)
		if diff := item.Rating.Wilson - want; diff > 1e-9 || diff < -1e-9 {
			t.Fatalf("expected wilson %v, got %v", want, item.Rating.Wilson)
		}
	})
}
