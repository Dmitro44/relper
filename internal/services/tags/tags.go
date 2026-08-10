package tags

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/tag"
	apperr "github.com/Dmitro44/relper/internal/errors"
	tagsrepo "github.com/Dmitro44/relper/internal/repositories/tags"
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

func (t *TagService) GetByID(ctx context.Context, id uuid.UUID) (*tag.Tag, error) {
	tagEntity, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tagsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tag not found")
		}
		return nil, apperr.NewInternalError("failed to get tag", err)
	}

	return tagEntity, nil
}

func (t *TagService) Create(ctx context.Context, input CreateInput) (*tag.Tag, error) {
	tagEntity, err := t.repo.Create(ctx, tagsrepo.CreateInput{
		Name: input.Name,
		Slug: input.Slug,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create tag", err)
	}
	slog.InfoContext(ctx, "tag created", "tag_id", tagEntity.ID)

	return tagEntity, nil
}

func (t *TagService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*tag.Tag, error) {
	tagEntity, err := t.repo.Update(ctx, id, tagsrepo.UpdateInput{
		Name: input.Name,
		Slug: input.Slug,
	})
	if err != nil {
		if errors.Is(err, tagsrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("tag not found")
		}
		return nil, apperr.NewInternalError("failed to update tag", err)
	}
	slog.InfoContext(ctx, "tag updated", "tag_id", tagEntity.ID)

	return tagEntity, nil
}

func (t *TagService) Delete(ctx context.Context, id uuid.UUID) error {
	err := t.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, tagsrepo.ErrNotFound) {
			return apperr.NewNotFoundError("tag not found")
		}
		return apperr.NewInternalError("failed to delete tag", err)
	}
	slog.InfoContext(ctx, "tag deleted", "tag_id", id)

	return nil
}
