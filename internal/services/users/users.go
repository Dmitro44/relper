package userservice

import (
	"context"

	userentity "github.com/Dmitro44/relper/internal/entities/user"
	usersrepo "github.com/Dmitro44/relper/internal/repositories/users"
	"github.com/Dmitro44/relper/internal/repositories/users/model"
	"github.com/google/uuid"
)

type UserService struct {
	repo usersrepo.Repository
}

func NewUserService(repo usersrepo.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) modelToEntity(m *model.User) *userentity.User {
	return &userentity.User{
		ID:       m.ID,
		Email:    m.Email,
		Nickname: m.Nickname,
		Habits:   m.Habits,
		Contacts: userentity.Contacts(m.Contacts),
	}
}

func (u *UserService) GetByID(ctx context.Context, id uuid.UUID) (*userentity.User, error) {
	userModel, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.modelToEntity(userModel), nil
}

type CreateInput struct {
	Email    string
	Nickname string
	Habits   []string
	Contacts userentity.Contacts
}

type UpdateInput struct {
	Email    *string
	Nickname *string
	Habits   *[]string
	Contacts *userentity.Contacts
}

func (u *UserService) Create(ctx context.Context, input CreateInput) (*userentity.User, error) {
	userModel, err := u.repo.Create(ctx, usersrepo.CreateInput{
		Email:    input.Email,
		Nickname: input.Nickname,
		Habits:   input.Habits,
		Contacts: input.Contacts,
	})
	if err != nil {
		return nil, err
	}

	return u.modelToEntity(userModel), nil
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*userentity.User, error) {
	userModel, err := u.repo.Update(ctx, id, usersrepo.UpdateInput{
		Email:    input.Email,
		Nickname: input.Nickname,
		Habits:   input.Habits,
		Contacts: input.Contacts,
	})
	if err != nil {
		return nil, err
	}

	return u.modelToEntity(userModel), nil
}
