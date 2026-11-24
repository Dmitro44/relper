package comments

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error)
	Create(ctx context.Context, input CreateInput) (*comment.Comment, error)
	Update(ctx context.Context, id, userID uuid.UUID, admin bool, input UpdateInput) (*comment.Comment, error)
	Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error
	ListByPlace(ctx context.Context, placeID uuid.UUID) ([]comment.WithAuthor, error)
	ListByTip(ctx context.Context, tipID uuid.UUID) ([]comment.WithAuthor, error)
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
