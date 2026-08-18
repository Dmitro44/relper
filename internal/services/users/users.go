package users

import (
	"context"
	"errors"
	"log/slog"

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
	Email        string
	PasswordHash []byte
	Nickname     string
	Habits       []string
	Contacts     userentity.Contacts
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

func (u *UserService) GetByEmail(ctx context.Context, email string) (*userentity.User, error) {
	userEntity, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, usersrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("user by email not found ")
		}
		return nil, err
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
	slog.InfoContext(ctx, "user created", "user_id", userEntity.ID)

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
	slog.InfoContext(ctx, "user updated", "user_id", userEntity.ID)

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
	slog.InfoContext(ctx, "user deleted", "user_id", id)

	return nil
}
