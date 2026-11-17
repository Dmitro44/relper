package likes

import (
	"context"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/Dmitro44/relper/internal/entities/like"
	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	likesrepo "github.com/Dmitro44/relper/internal/repositories/likes"
	"github.com/google/uuid"
)

type TipGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error)
}

type CommentGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error)
}

type LikeService struct {
	repo     likesrepo.Repository
	tips     TipGetter
	comments CommentGetter
}

func NewLikeService(repo likesrepo.Repository, tips TipGetter, comments CommentGetter) *LikeService {
	return &LikeService{
		repo:     repo,
		tips:     tips,
		comments: comments,
	}
}

func (s *LikeService) LikeTip(ctx context.Context, userID, tipID uuid.UUID) (*like.LikeSummary, error) {
	tipEntity, err := s.tips.GetByID(ctx, tipID)
	if err != nil {
		return nil, err
	}
	if tipEntity.UserID == userID {
		return nil, apperr.NewForbiddenError("cannot like your own tip")
	}

	if err := s.repo.LikeTip(ctx, userID, tipID); err != nil {
		return nil, apperr.NewInternalError("failed to like tip", err)
	}
	slog.InfoContext(ctx, "tip liked", "tip_id", tipID, "user_id", userID)

	return s.TipSummary(ctx, tipID, &userID)
}

func (s *LikeService) UnlikeTip(ctx context.Context, userID, tipID uuid.UUID) (*like.LikeSummary, error) {
	if err := s.repo.UnlikeTip(ctx, userID, tipID); err != nil {
		return nil, apperr.NewInternalError("failed to unlike tip", err)
	}
	slog.InfoContext(ctx, "tip unliked", "tip_id", tipID, "user_id", userID)

	return s.TipSummary(ctx, tipID, &userID)
}

func (s *LikeService) LikeComment(ctx context.Context, userID, commentID uuid.UUID) (*like.LikeSummary, error) {
	commentEntity, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if commentEntity.UserID == userID {
		return nil, apperr.NewForbiddenError("cannot like your own comment")
	}

	if err := s.repo.LikeComment(ctx, userID, commentID); err != nil {
		return nil, apperr.NewInternalError("failed to like comment", err)
	}
	slog.InfoContext(ctx, "comment liked", "comment_id", commentID, "user_id", userID)

	return s.CommentSummary(ctx, commentID, &userID)
}

func (s *LikeService) UnlikeComment(ctx context.Context, userID, commentID uuid.UUID) (*like.LikeSummary, error) {
	if err := s.repo.UnlikeComment(ctx, userID, commentID); err != nil {
		return nil, apperr.NewInternalError("failed to unlike comment", err)
	}
	slog.InfoContext(ctx, "comment unliked", "comment_id", commentID, "user_id", userID)

	return s.CommentSummary(ctx, commentID, &userID)
}

func (s *LikeService) TipSummary(ctx context.Context, tipID uuid.UUID, viewer *uuid.UUID) (*like.LikeSummary, error) {
	if _, err := s.tips.GetByID(ctx, tipID); err != nil {
		return nil, err
	}

	count, err := s.repo.CountTip(ctx, tipID)
	if err != nil {
		return nil, apperr.NewInternalError("failed to count tip likes", err)
	}

	summary := like.LikeSummary{Count: count}
	if viewer != nil {
		liked, err := s.repo.LikedTip(ctx, *viewer, tipID)
		if err != nil {
			return nil, apperr.NewInternalError("failed to get tip like", err)
		}
		summary.MyLike = &liked
	}

	return &summary, nil
}

func (s *LikeService) CommentSummary(ctx context.Context, commentID uuid.UUID, viewer *uuid.UUID) (*like.LikeSummary, error) {
	if _, err := s.comments.GetByID(ctx, commentID); err != nil {
		return nil, err
	}

	count, err := s.repo.CountComment(ctx, commentID)
	if err != nil {
		return nil, apperr.NewInternalError("failed to count comment likes", err)
	}

	summary := like.LikeSummary{Count: count}
	if viewer != nil {
		liked, err := s.repo.LikedComment(ctx, *viewer, commentID)
		if err != nil {
			return nil, apperr.NewInternalError("failed to get comment like", err)
		}
		summary.MyLike = &liked
	}

	return &summary, nil
}
