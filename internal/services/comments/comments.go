package comments

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/Dmitro44/relper/internal/entities/place"
	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	commentsrepo "github.com/Dmitro44/relper/internal/repositories/comments"
	"github.com/google/uuid"
)

type PlaceGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error)
}

type TipGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error)
}

type CommentService struct {
	repo   commentsrepo.Repository
	places PlaceGetter
	tips   TipGetter
}

type CreateInput struct {
	PlaceID *uuid.UUID
	TipID   *uuid.UUID
	UserID  uuid.UUID
	Content string
}

type UpdateInput struct {
	Content *string
}

func NewCommmentService(repo commentsrepo.Repository, places PlaceGetter, tips TipGetter) *CommentService {
	return &CommentService{
		repo:   repo,
		places: places,
		tips:   tips,
	}
}

func (c *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	commentEntity, err := c.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, commentsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("comment not found")
		}
		return nil, apperr.NewInternalError("failed to get comment", err)
	}

	return commentEntity, nil
}

func (c *CommentService) ListByPlace(ctx context.Context, placeID uuid.UUID) ([]comment.WithAuthor, error) {
	if _, err := c.places.GetByID(ctx, placeID); err != nil {
		return nil, err
	}

	items, err := c.repo.ListByPlace(ctx, placeID)
	if err != nil {
		return nil, apperr.NewInternalError("failed to list place comments", err)
	}

	return items, nil
}

func (c *CommentService) ListByTip(ctx context.Context, tipID uuid.UUID) ([]comment.WithAuthor, error) {
	if _, err := c.tips.GetByID(ctx, tipID); err != nil {
		return nil, err
	}

	items, err := c.repo.ListByTip(ctx, tipID)
	if err != nil {
		return nil, apperr.NewInternalError("failed to list tip comments", err)
	}

	return items, nil
}

func (c *CommentService) Create(ctx context.Context, input CreateInput) (*comment.Comment, error) {
	if (input.PlaceID == nil) == (input.TipID == nil) {
		return nil, apperr.NewBadRequestError("exactly one of place_id or tip_id is required")
	}

	if input.PlaceID != nil {
		if _, err := c.places.GetByID(ctx, *input.PlaceID); err != nil {
			return nil, err
		}
	}
	if input.TipID != nil {
		if _, err := c.tips.GetByID(ctx, *input.TipID); err != nil {
			return nil, err
		}
	}

	commentEntity, err := c.repo.Create(ctx, commentsrepo.CreateInput{
		PlaceID: input.PlaceID,
		TipID:   input.TipID,
		UserID:  input.UserID,
		Content: input.Content,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create comment", err)
	}
	slog.InfoContext(ctx, "comment created", "comment_id", commentEntity.ID)

	return commentEntity, nil
}

func (c *CommentService) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input UpdateInput) (*comment.Comment, error) {
	if input.Content == nil || *input.Content == "" {
		return nil, apperr.NewBadRequestError("content is required")
	}

	commentEntity, err := c.repo.Update(ctx, id, userID, admin, commentsrepo.UpdateInput{
		Content: input.Content,
	})
	if err != nil {
		if errors.Is(err, commentsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("comment not found")
		}
		return nil, apperr.NewInternalError("failed to update comment", err)
	}
	slog.InfoContext(ctx, "comment updated", "comment_id", commentEntity.ID)

	return commentEntity, nil
}

func (c *CommentService) Delete(ctx context.Context, id, callerID uuid.UUID, admin bool) error {
	err := c.repo.Delete(ctx, id, callerID, admin)
	if err != nil {
		if errors.Is(err, commentsrepo.ErrNotFound) {
			return apperr.NewNotFoundError("comment not found")
		}
		return apperr.NewInternalError("failed to delete comment", err)
	}
	slog.InfoContext(ctx, "comment deleted", "comment_id", id)

	return nil
}
