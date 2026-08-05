package comments

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error)
	Create(ctx context.Context, input CreateInput) (*comment.Comment, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*comment.Comment, error)
}

type CreateInput struct {
	TargetType string
	TargetID   uuid.UUID
	UserID     uuid.UUID
	Content    string
}

type UpdateInput struct {
	TargetType *string
	TargetID   *uuid.UUID
	UserID     *uuid.UUID
	Content    *string
}
