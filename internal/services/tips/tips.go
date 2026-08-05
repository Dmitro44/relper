package tips

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/tip"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips"
	"github.com/Dmitro44/relper/internal/repositories/tips/model"
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

func (t *TipService) modelToEntity(m *model.Tip) *tip.Tip {
	return &tip.Tip{
		ID:         m.ID,
		LocationID: m.LocationID,
		UserID:     m.UserID,
		Category:   m.Category,
		Title:      m.Title,
		Content:    m.Content,
		Upvotes:    m.Upvotes,
	}
}

func (t *TipService) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	tipModel, err := t.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return t.modelToEntity(tipModel), nil
}

func (t *TipService) Create(ctx context.Context, input CreateInput) (*tip.Tip, error) {
	tipModel, err := t.repo.Create(ctx, tipsrepo.CreateInput{
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

	return t.modelToEntity(tipModel), nil
}

func (t *TipService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*tip.Tip, error) {
	tipModel, err := t.repo.Update(ctx, id, tipsrepo.UpdateInput{
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

	return t.modelToEntity(tipModel), nil
}
