package users

import (
	"context"
	"errors"

	userentity "github.com/Dmitro44/relper/internal/entities/user"
	apperr "github.com/Dmitro44/relper/internal/errors"
	usersrepo "github.com/Dmitro44/relper/internal/repositories/users"
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

func (u *UserService) GetByID(ctx context.Context, id uuid.UUID) (*userentity.User, error) {
	userEntity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, usersrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("user not found")
		}
		return nil, apperr.NewInternalError("failed to get user", err)
	}

	return userEntity, nil
}

func (u *UserService) Create(ctx context.Context, input CreateInput) (*userentity.User, error) {
	userEntity, err := u.repo.Create(ctx, usersrepo.CreateInput{
		Email:    input.Email,
		Nickname: input.Nickname,
		Habits:   input.Habits,
		Contacts: input.Contacts,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create user", err)
	}

	return userEntity, nil
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*userentity.User, error) {
	userEntity, err := u.repo.Update(ctx, id, usersrepo.UpdateInput{
		Email:    input.Email,
		Nickname: input.Nickname,
		Habits:   input.Habits,
		Contacts: input.Contacts,
	})
	if err != nil {
		if errors.Is(err, usersrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("user not found")
		}
		return nil, apperr.NewInternalError("failed to update user", err)
	}

	return userEntity, nil
}

func (u *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, usersrepo.ErrNotFound) {
			return apperr.NewNotFoundError("user not found")
		}
		return apperr.NewInternalError("failed to delete user", err)
	}

	return nil
}
