package tips

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/tip"
	apperr "github.com/Dmitro44/relper/internal/errors"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips"
	"github.com/google/uuid"
)

type TipService struct {
	repo tipsrepo.Repository
}

type CreateInput struct {
	LocationID uuid.UUID
	UserID     uuid.UUID
	Category   string
	Title      string
	Content    string
	Upvotes    int32
}

type UpdateInput struct {
	LocationID *uuid.UUID
	Category   *string
	Title      *string
	Content    *string
	Upvotes    *int32
}

func NewTipService(repo tipsrepo.Repository) *TipService {
	return &TipService{
		repo: repo,
	}
}

func (t *TipService) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	tipEntity, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tip not found")
		}
		return nil, apperr.NewInternalError("failed to get tip", err)
	}

	return tipEntity, nil
}

func (t *TipService) Create(ctx context.Context, input CreateInput) (*tip.Tip, error) {
	tipEntity, err := t.repo.Create(ctx, tipsrepo.CreateInput{
		LocationID: input.LocationID,
		UserID:     input.UserID,
		Category:   input.Category,
		Title:      input.Title,
		Content:    input.Content,
		Upvotes:    input.Upvotes,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create tip", err)
	}
	slog.InfoContext(ctx, "tip created", "tip_id", tipEntity.ID)

	return tipEntity, nil
}

func (t *TipService) Update(ctx context.Context, id, userID uuid.UUID, input UpdateInput) (*tip.Tip, error) {
	tipEntity, err := t.repo.Update(ctx, id, userID, tipsrepo.UpdateInput{
		LocationID: input.LocationID,
		Category:   input.Category,
		Title:      input.Title,
		Content:    input.Content,
		Upvotes:    input.Upvotes,
	})
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tip not found")
		}
		return nil, apperr.NewInternalError("failed to update tip", err)
	}
	slog.InfoContext(ctx, "tip updated", "tip_id", tipEntity.ID)

	return tipEntity, nil
}

func (t *TipService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	err := t.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return apperr.NewNotFoundError("tip not found")
		}
		return apperr.NewInternalError("failed to delete tip", err)
	}
	slog.InfoContext(ctx, "tip deleted", "tip_id", id)

	return nil
}

func (t *TipService) Upvote(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	updTip, err := t.repo.Upvote(ctx, id)
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tip for upvote not found")
		}
		return nil, apperr.NewInternalError("failed to upvote tip", err)
	}
	slog.InfoContext(ctx, "tip upvoted", "tip_id", updTip.ID)

	return updTip, nil
}

func (t *TipService) RemoveUpvote(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	updTip, err := t.repo.RemoveUpvote(ctx, id)
	if err != nil {
		if errors.Is(err, tipsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tip for remove upvote not found")
		}
		return nil, apperr.NewInternalError("failed to remove upvote from tip", err)
	}
	slog.InfoContext(ctx, "tip removed upvote", "tip_id", updTip.ID)

	return updTip, nil
}
