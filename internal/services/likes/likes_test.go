package likes

import (
	"context"
	"errors"
	"testing"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/google/uuid"
)

type call struct {
	userID   uuid.UUID
	targetID uuid.UUID
}

type fakeRepo struct {
	likeTipCalls       []call
	likeCommentCalls   []call
	unlikeTipCalls     []call
	unlikeCommentCalls []call

	likeErr   error
	unlikeErr error

	likedTip        bool
	likedTipErr     error
	likedComment    bool
	likedCommentErr error
	countTip        int64
	countComment    int64
	countErr        error
}

func (f *fakeRepo) LikeTip(ctx context.Context, userID, tipID uuid.UUID) error {
	f.likeTipCalls = append(f.likeTipCalls, call{userID, tipID})
	return f.likeErr
}

func (f *fakeRepo) LikeComment(ctx context.Context, userID, commentID uuid.UUID) error {
	f.likeCommentCalls = append(f.likeCommentCalls, call{userID, commentID})
	return f.likeErr
}

func (f *fakeRepo) UnlikeTip(ctx context.Context, userID, tipID uuid.UUID) error {
	f.unlikeTipCalls = append(f.unlikeTipCalls, call{userID, tipID})
	return f.unlikeErr
}

func (f *fakeRepo) UnlikeComment(ctx context.Context, userID, commentID uuid.UUID) error {
	f.unlikeCommentCalls = append(f.unlikeCommentCalls, call{userID, commentID})
	return f.unlikeErr
}

func (f *fakeRepo) LikedTip(ctx context.Context, userID, tipID uuid.UUID) (bool, error) {
	return f.likedTip, f.likedTipErr
}

func (f *fakeRepo) LikedComment(ctx context.Context, userID, commentID uuid.UUID) (bool, error) {
	return f.likedComment, f.likedCommentErr
}

func (f *fakeRepo) CountTip(ctx context.Context, tipID uuid.UUID) (int64, error) {
	return f.countTip, f.countErr
}

func (f *fakeRepo) CountComment(ctx context.Context, commentID uuid.UUID) (int64, error) {
	return f.countComment, f.countErr
}

type fakeTipGetter struct {
	tip tip.Tip
	err error
}

func (f *fakeTipGetter) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &f.tip, nil
}

type fakeCommentGetter struct {
	comment comment.Comment
	err     error
}

func (f *fakeCommentGetter) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &f.comment, nil
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

func TestLikeTip(t *testing.T) {
	userID, tipID := uuid.New(), uuid.New()

	t.Run("tip not found", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewLikeService(repo, &fakeTipGetter{err: apperr.NewNotFoundError("tip not found")}, &fakeCommentGetter{})

		_, err := svc.LikeTip(context.Background(), userID, tipID)

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
		if len(repo.likeTipCalls) != 0 {
			t.Fatal("expected no like")
		}
	})

	t.Run("self like forbidden", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewLikeService(repo, &fakeTipGetter{tip: tip.Tip{ID: tipID, UserID: userID}}, &fakeCommentGetter{})

		_, err := svc.LikeTip(context.Background(), userID, tipID)

		assertErrorCode(t, err, apperr.ErrorCodeForbidden)
		if len(repo.likeTipCalls) != 0 {
			t.Fatal("expected no like")
		}
	})

	t.Run("likes tip", func(t *testing.T) {
		repo := &fakeRepo{countTip: 3, likedTip: true}
		svc := NewLikeService(repo, &fakeTipGetter{tip: tip.Tip{ID: tipID, UserID: uuid.New()}}, &fakeCommentGetter{})

		summary, err := svc.LikeTip(context.Background(), userID, tipID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.likeTipCalls) != 1 || repo.likeTipCalls[0] != (call{userID, tipID}) {
			t.Fatalf("unexpected like calls: %+v", repo.likeTipCalls)
		}
		if summary.Count != 3 {
			t.Fatalf("expected count 3, got %d", summary.Count)
		}
		if summary.MyLike == nil || !*summary.MyLike {
			t.Fatalf("expected my_like true, got %v", summary.MyLike)
		}
	})

	t.Run("like failure", func(t *testing.T) {
		repo := &fakeRepo{likeErr: errors.New("boom")}
		svc := NewLikeService(repo, &fakeTipGetter{tip: tip.Tip{ID: tipID, UserID: uuid.New()}}, &fakeCommentGetter{})

		_, err := svc.LikeTip(context.Background(), userID, tipID)

		assertErrorCode(t, err, apperr.ErrorCodeInternal)
	})
}

func TestUnlikeTip(t *testing.T) {
	userID, tipID := uuid.New(), uuid.New()

	t.Run("idempotent", func(t *testing.T) {
		repo := &fakeRepo{countTip: 0}
		svc := NewLikeService(repo, &fakeTipGetter{tip: tip.Tip{ID: tipID}}, &fakeCommentGetter{})

		for range 2 {
			summary, err := svc.UnlikeTip(context.Background(), userID, tipID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if summary.Count != 0 {
				t.Fatalf("expected count 0, got %d", summary.Count)
			}
			if summary.MyLike == nil || *summary.MyLike {
				t.Fatalf("expected my_like false, got %v", summary.MyLike)
			}
		}
		if len(repo.unlikeTipCalls) != 2 || repo.unlikeTipCalls[0] != (call{userID, tipID}) {
			t.Fatalf("unexpected unlike calls: %+v", repo.unlikeTipCalls)
		}
	})

	t.Run("unlike failure", func(t *testing.T) {
		repo := &fakeRepo{unlikeErr: errors.New("boom")}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{})

		_, err := svc.UnlikeTip(context.Background(), userID, tipID)

		assertErrorCode(t, err, apperr.ErrorCodeInternal)
	})
}

func TestLikeComment(t *testing.T) {
	userID, commentID := uuid.New(), uuid.New()

	t.Run("comment not found", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{err: apperr.NewNotFoundError("comment not found")})

		_, err := svc.LikeComment(context.Background(), userID, commentID)

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
		if len(repo.likeCommentCalls) != 0 {
			t.Fatal("expected no like")
		}
	})

	t.Run("self like forbidden", func(t *testing.T) {
		repo := &fakeRepo{}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{comment: comment.Comment{ID: commentID, UserID: userID}})

		_, err := svc.LikeComment(context.Background(), userID, commentID)

		assertErrorCode(t, err, apperr.ErrorCodeForbidden)
		if len(repo.likeCommentCalls) != 0 {
			t.Fatal("expected no like")
		}
	})

	t.Run("likes comment", func(t *testing.T) {
		repo := &fakeRepo{countComment: 2, likedComment: true}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{comment: comment.Comment{ID: commentID, UserID: uuid.New()}})

		summary, err := svc.LikeComment(context.Background(), userID, commentID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.likeCommentCalls) != 1 || repo.likeCommentCalls[0] != (call{userID, commentID}) {
			t.Fatalf("unexpected like calls: %+v", repo.likeCommentCalls)
		}
		if summary.Count != 2 {
			t.Fatalf("expected count 2, got %d", summary.Count)
		}
		if summary.MyLike == nil || !*summary.MyLike {
			t.Fatalf("expected my_like true, got %v", summary.MyLike)
		}
	})
}

func TestUnlikeComment(t *testing.T) {
	userID, commentID := uuid.New(), uuid.New()

	repo := &fakeRepo{countComment: 0, likedComment: false}
	svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{comment: comment.Comment{ID: commentID}})

	summary, err := svc.UnlikeComment(context.Background(), userID, commentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.unlikeCommentCalls) != 1 || repo.unlikeCommentCalls[0] != (call{userID, commentID}) {
		t.Fatalf("unexpected unlike calls: %+v", repo.unlikeCommentCalls)
	}
	if summary.MyLike == nil || *summary.MyLike {
		t.Fatalf("expected my_like false, got %v", summary.MyLike)
	}
}

func TestSummary(t *testing.T) {
	t.Run("anonymous tip summary", func(t *testing.T) {
		repo := &fakeRepo{countTip: 7}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{})

		summary, err := svc.TipSummary(context.Background(), uuid.New(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if summary.Count != 7 || summary.MyLike != nil {
			t.Fatalf("unexpected summary: %+v", summary)
		}
	})

	t.Run("tip not found", func(t *testing.T) {
		svc := NewLikeService(&fakeRepo{}, &fakeTipGetter{err: apperr.NewNotFoundError("tip not found")}, &fakeCommentGetter{})

		_, err := svc.TipSummary(context.Background(), uuid.New(), nil)

		assertErrorCode(t, err, apperr.ErrorCodeNotFound)
	})

	t.Run("comment summary", func(t *testing.T) {
		commentID, viewer := uuid.New(), uuid.New()
		repo := &fakeRepo{countComment: 4, likedComment: true}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{comment: comment.Comment{ID: commentID}})

		summary, err := svc.CommentSummary(context.Background(), commentID, &viewer)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if summary.Count != 4 || summary.MyLike == nil || !*summary.MyLike {
			t.Fatalf("unexpected summary: %+v", summary)
		}
	})

	t.Run("count failure", func(t *testing.T) {
		repo := &fakeRepo{countErr: errors.New("boom")}
		svc := NewLikeService(repo, &fakeTipGetter{}, &fakeCommentGetter{})

		_, err := svc.CommentSummary(context.Background(), uuid.New(), nil)

		assertErrorCode(t, err, apperr.ErrorCodeInternal)
	})
}
