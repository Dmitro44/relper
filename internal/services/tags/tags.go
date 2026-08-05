package tags

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/tag"
	tagsrepo "github.com/Dmitro44/relper/internal/repositories/tags"
	"github.com/Dmitro44/relper/internal/repositories/tags/model"
	"github.com/google/uuid"
)

type TagService struct {
	repo tagsrepo.Repository
}

type CreateInput struct {
	Name string
	Slug string
}

type UpdateInput struct {
	Name *string
	Slug *string
}

func NewTagService(repo tagsrepo.Repository) *TagService {
	return &TagService{
		repo: repo,
	}
}

func (t *TagService) modelToEntity(m *model.Tag) *tag.Tag {
	return &tag.Tag{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		CreatedAt: m.CreatedAt,
	}
}

func (t *TagService) GetByID(ctx context.Context, id uuid.UUID) (*tag.Tag, error) {
	tagModel, err := t.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return t.modelToEntity(tagModel), nil
}

func (t *TagService) Create(ctx context.Context, input CreateInput) (*tag.Tag, error) {
	tagModel, err := t.repo.Create(ctx, tagsrepo.CreateInput{
		Name: input.Name,
		Slug: input.Slug,
	})
	if err != nil {
		return nil, err
	}

	return t.modelToEntity(tagModel), nil
}

func (t *TagService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*tag.Tag, error) {
	tagModel, err := t.repo.Update(ctx, id, tagsrepo.UpdateInput{
		Name: input.Name,
		Slug: input.Slug,
	})
	if err != nil {
		return nil, err
	}

	return t.modelToEntity(tagModel), nil
}
