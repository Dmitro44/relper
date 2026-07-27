package users

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/users/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, input CreateInput) (*model.User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*model.User, error)
}

type CreateInput struct {
	Email    string
	Nickname string
	Habits   []string
	Contacts model.Contacts
}

type UpdateInput struct {
	Email    *string
	Nickname *string
	Habits   *[]string
	Contacts *model.Contacts
}
