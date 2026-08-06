package tips

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/tip"
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
	UserID     *uuid.UUID
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
		return nil, err
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
		return nil, err
	}

	return tipEntity, nil
}

func (t *TipService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*tip.Tip, error) {
	tipEntity, err := t.repo.Update(ctx, id, tipsrepo.UpdateInput{
		LocationID: input.LocationID,
		UserID:     input.UserID,
		Category:   input.Category,
		Title:      input.Title,
		Content:    input.Content,
		Upvotes:    input.Upvotes,
	})
	if err != nil {
		return nil, err
	}

	return tipEntity, nil
}

func (t *TipService) Delete(ctx context.Context, id uuid.UUID) error {
	err := t.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TipService) Upvote(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	updTip, err := t.repo.Upvote(ctx, id)
	if err != nil {
		return nil, err
	}

	return updTip, nil
}

func (t *TipService) RemoveUpvote(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	updTip, err := t.repo.RemoveUpvote(ctx, id)
	if err != nil {
		return nil, err
	}

	return updTip, nil
}
