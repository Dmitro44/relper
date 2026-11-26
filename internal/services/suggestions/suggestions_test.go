package suggestions

import (
	"context"
	"errors"
	"testing"

	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/entities/suggestion"
	apperr "github.com/Dmitro44/relper/internal/errors"
	suggestionsrepo "github.com/Dmitro44/relper/internal/repositories/suggestions"
	"github.com/google/uuid"
)

type fakeRepository struct {
	created    *suggestion.Suggestion
	createErr  error
	createIn   suggestionsrepo.CreateInput
	createCall bool

	listMine     []suggestion.Suggestion
	listMineErr  error
	listMineUser uuid.UUID

	list      []suggestion.Suggestion
	listErr   error
	listState string

	approve      *suggestion.Suggestion
	approveErr   error
	approveID    uuid.UUID
	approveUsers []uuid.UUID

	reject    *suggestion.Suggestion
	rejectErr error
	rejectWhy string
}

func (f *fakeRepository) Create(ctx context.Context, input suggestionsrepo.CreateInput) (*suggestion.Suggestion, error) {
	f.createCall = true
	f.createIn = input
	return f.created, f.createErr
}

func (f *fakeRepository) ListMine(ctx context.Context, userID uuid.UUID) ([]suggestion.Suggestion, error) {
	f.listMineUser = userID
	return f.listMine, f.listMineErr
}

func (f *fakeRepository) ListByStatus(ctx context.Context, status string) ([]suggestion.Suggestion, error) {
	f.listState = status
	return f.list, f.listErr
}

func (f *fakeRepository) Approve(ctx context.Context, id, reviewerID uuid.UUID) (*suggestion.Suggestion, error) {
	f.approveID = id
	f.approveUsers = []uuid.UUID{reviewerID}
	return f.approve, f.approveErr
}

func (f *fakeRepository) Reject(ctx context.Context, id, reviewerID uuid.UUID, reason string) (*suggestion.Suggestion, error) {
	f.rejectWhy = reason
	return f.reject, f.rejectErr
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

func assertCode(t *testing.T, err error, want apperr.ErrorCode) {
	t.Helper()

	var appErr *apperr.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *apperr.AppError, got %v", err)
	}
	if appErr.Code != want {
		t.Fatalf("expected code %s, got %s", want, appErr.Code)
	}
}

func validInput() CreateInput {
	return CreateInput{
		UserID:     uuid.New(),
		LocationID: uuid.New(),
		Name:       "Кофейня",
		Category:   "food",
	}
}

func TestCreate(t *testing.T) {
	t.Run("unknown category", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewSuggestionService(repo, &fakeLocationGetter{})

		_, err := svc.Create(context.Background(), CreateInput{
			LocationID: uuid.New(),
			Name:       "X",
			Category:   "casino",
		})

		assertCode(t, err, apperr.ErrorCodeBadRequest)
		if repo.createCall {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("blank name", func(t *testing.T) {
		svc := NewSuggestionService(&fakeRepository{}, &fakeLocationGetter{})

		input := validInput()
		input.Name = "   "
		_, err := svc.Create(context.Background(), input)

		assertCode(t, err, apperr.ErrorCodeBadRequest)
	})

	t.Run("location not found", func(t *testing.T) {
		repo := &fakeRepository{}
		locations := &fakeLocationGetter{err: apperr.NewNotFoundError("location not found")}
		svc := NewSuggestionService(repo, locations)

		_, err := svc.Create(context.Background(), validInput())

		assertCode(t, err, apperr.ErrorCodeNotFound)
		if repo.createCall {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("trims category and passes user", func(t *testing.T) {
		repo := &fakeRepository{created: &suggestion.Suggestion{ID: uuid.New()}}
		svc := NewSuggestionService(repo, &fakeLocationGetter{})

		input := validInput()
		input.Category = " food "
		_, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.createIn.Category != "food" {
			t.Fatalf("expected trimmed category, got %q", repo.createIn.Category)
		}
		if repo.createIn.UserID != input.UserID {
			t.Fatal("expected user id forwarded")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("defaults to pending", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewSuggestionService(repo, &fakeLocationGetter{})

		if _, err := svc.List(context.Background(), ""); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.listState != suggestion.StatusPending {
			t.Fatalf("expected pending, got %q", repo.listState)
		}
	})

	t.Run("unknown status", func(t *testing.T) {
		svc := NewSuggestionService(&fakeRepository{}, &fakeLocationGetter{})

		_, err := svc.List(context.Background(), "archived")

		assertCode(t, err, apperr.ErrorCodeBadRequest)
	})
}

func TestApprove(t *testing.T) {
	t.Run("not pending", func(t *testing.T) {
		svc := NewSuggestionService(&fakeRepository{approveErr: suggestionsrepo.ErrNotPending}, &fakeLocationGetter{})

		_, err := svc.Approve(context.Background(), uuid.New(), uuid.New())

		assertCode(t, err, apperr.ErrorCodeConflict)
	})

	t.Run("not found", func(t *testing.T) {
		svc := NewSuggestionService(&fakeRepository{approveErr: suggestionsrepo.ErrNotFound}, &fakeLocationGetter{})

		_, err := svc.Approve(context.Background(), uuid.New(), uuid.New())

		assertCode(t, err, apperr.ErrorCodeNotFound)
	})
}

func TestReject(t *testing.T) {
	t.Run("reason required", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewSuggestionService(repo, &fakeLocationGetter{})

		_, err := svc.Reject(context.Background(), uuid.New(), uuid.New(), "  ")

		assertCode(t, err, apperr.ErrorCodeBadRequest)
		if repo.rejectWhy != "" {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("not pending", func(t *testing.T) {
		svc := NewSuggestionService(&fakeRepository{rejectErr: suggestionsrepo.ErrNotPending}, &fakeLocationGetter{})

		_, err := svc.Reject(context.Background(), uuid.New(), uuid.New(), "duplicate")

		assertCode(t, err, apperr.ErrorCodeConflict)
	})
}
