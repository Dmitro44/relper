package comments

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/comments/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	Create(ctx context.Context, input CreateInput) (*model.Comment, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*model.Comment, error)
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
