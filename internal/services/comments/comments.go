package comments

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/comment"
	commentsrepo "github.com/Dmitro44/relper/internal/repositories/comments"
	"github.com/Dmitro44/relper/internal/repositories/comments/model"
	"github.com/google/uuid"
)

type CommentService struct {
	repo commentsrepo.Repository
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

func NewCommmentService(repo commentsrepo.Repository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (c *CommentService) modelToEntity(m *model.Comment) *comment.Comment {
	return &comment.Comment{
		ID:         m.ID,
		TargetType: m.TargetType,
		TargetID:   m.TargetID,
		UserID:     m.UserID,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
	}
}

func (c *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	commentModel, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return c.modelToEntity(commentModel), nil
}

func (c *CommentService) Create(ctx context.Context, input CreateInput) (*comment.Comment, error) {
	// checks for ids (target, user)
	//
	// another business logic
	// if input.TargetType == "tip" {
	// 	_, err := c.repo.
	// }

	commentModel, err := c.repo.Create(ctx, commentsrepo.CreateInput{
		TargetType: input.TargetType,
		TargetID:   input.TargetID,
		UserID:     input.UserID,
		Content:    input.Content,
	})
	if err != nil {
		return nil, err
	}

	return c.modelToEntity(commentModel), nil
}

func (c *CommentService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*comment.Comment, error) {
	// checks for ids (target, user)
	//
	// another business logic

	commentModel, err := c.repo.Update(ctx, id, commentsrepo.UpdateInput{
		TargetType: input.TargetType,
		TargetID:   input.TargetID,
		UserID:     input.UserID,
		Content:    input.Content,
	})
	if err != nil {
		return nil, err
	}

	return c.modelToEntity(commentModel), nil
}
