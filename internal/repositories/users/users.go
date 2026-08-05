package users

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	Create(ctx context.Context, input CreateInput) (*user.User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*user.User, error)
}

type CreateInput struct {
	Email    string
	Nickname string
	Habits   []string
	Contacts user.Contacts
}

type UpdateInput struct {
	Email    *string
	Nickname *string
	Habits   *[]string
	Contacts *user.Contacts
}
