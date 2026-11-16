package comments

import (
	"context"
	"errors"
	"testing"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	commentsrepo "github.com/Dmitro44/relper/internal/repositories/comments"
	"github.com/google/uuid"
)

type fakeRepo struct {
	created commentsrepo.CreateInput
	updated commentsrepo.UpdateInput
}

func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	return &comment.Comment{ID: id}, nil
}

func (f *fakeRepo) Create(ctx context.Context, input commentsrepo.CreateInput) (*comment.Comment, error) {
	f.created = input
	return &comment.Comment{
		ID:      uuid.New(),
		PlaceID: input.PlaceID,
		TipID:   input.TipID,
		UserID:  input.UserID,
		Content: input.Content,
	}, nil
}

func (f *fakeRepo) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input commentsrepo.UpdateInput) (*comment.Comment, error) {
	f.updated = input
	return &comment.Comment{ID: id, UserID: userID, Content: *input.Content}, nil
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

type fakeTipGetter struct {
	err error
}

func (f *fakeTipGetter) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &tip.Tip{ID: id}, nil
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
	content := "nice place"

	t.Run("both targets", func(t *testing.T) {
		placeID, tipID := uuid.New(), uuid.New()
		svc := NewCommmentService(&fakeRepo{}, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			PlaceID: &placeID,
			TipID:   &tipID,
			UserID:  userID,
			Content: content,
		})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("neither target", func(t *testing.T) {
		svc := NewCommmentService(&fakeRepo{}, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			UserID:  userID,
			Content: content,
		})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("place not found", func(t *testing.T) {
		placeID := uuid.New()
		svc := NewCommmentService(
			&fakeRepo{},
			&fakePlaceGetter{err: apperr.NewNotFoundError("place not found")},
			&fakeTipGetter{},
		)

		_, err := svc.Create(context.Background(), CreateInput{
			PlaceID: &placeID,
			UserID:  userID,
			Content: content,
		})

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})

	t.Run("tip not found", func(t *testing.T) {
		tipID := uuid.New()
		svc := NewCommmentService(
			&fakeRepo{},
			&fakePlaceGetter{},
			&fakeTipGetter{err: apperr.NewNotFoundError("tip not found")},
		)

		_, err := svc.Create(context.Background(), CreateInput{
			TipID:   &tipID,
			UserID:  userID,
			Content: content,
		})

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})

	t.Run("place comment", func(t *testing.T) {
		placeID := uuid.New()
		repo := &fakeRepo{}
		svc := NewCommmentService(repo, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			PlaceID: &placeID,
			UserID:  userID,
			Content: content,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.created.PlaceID == nil || *repo.created.PlaceID != placeID {
			t.Fatalf("expected place_id %s, got %v", placeID, repo.created.PlaceID)
		}
		if repo.created.TipID != nil {
			t.Fatalf("expected nil tip_id, got %s", *repo.created.TipID)
		}
		if repo.created.UserID != userID || repo.created.Content != content {
			t.Fatalf("unexpected create input: %+v", repo.created)
		}
	})

	t.Run("tip comment", func(t *testing.T) {
		tipID := uuid.New()
		repo := &fakeRepo{}
		svc := NewCommmentService(repo, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			TipID:   &tipID,
			UserID:  userID,
			Content: content,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.created.TipID == nil || *repo.created.TipID != tipID {
			t.Fatalf("expected tip_id %s, got %v", tipID, repo.created.TipID)
		}
		if repo.created.PlaceID != nil {
			t.Fatalf("expected nil place_id, got %s", *repo.created.PlaceID)
		}
		if repo.created.UserID != userID || repo.created.Content != content {
			t.Fatalf("unexpected create input: %+v", repo.created)
		}
	})
}

func TestUpdate(t *testing.T) {
	empty := ""

	t.Run("nil content", func(t *testing.T) {
		svc := NewCommmentService(&fakeRepo{}, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Update(context.Background(), uuid.New(), uuid.New(), false, UpdateInput{})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("empty content", func(t *testing.T) {
		svc := NewCommmentService(&fakeRepo{}, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Update(context.Background(), uuid.New(), uuid.New(), false, UpdateInput{
			Content: &empty,
		})

		assertErrorCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("updates content", func(t *testing.T) {
		content := "edited"
		repo := &fakeRepo{}
		svc := NewCommmentService(repo, &fakePlaceGetter{}, &fakeTipGetter{})

		_, err := svc.Update(context.Background(), uuid.New(), uuid.New(), false, UpdateInput{
			Content: &content,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.updated.Content == nil || *repo.updated.Content != content {
			t.Fatalf("expected content %q, got %v", content, repo.updated.Content)
		}
	})
}
