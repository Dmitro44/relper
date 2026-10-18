package users

import (
	"context"
	"errors"
	"log/slog"

	userentity "github.com/Dmitro44/relper/internal/entities/user"
	apperr "github.com/Dmitro44/relper/internal/errors"
	rolesrepo "github.com/Dmitro44/relper/internal/repositories/roles"
	usersrepo "github.com/Dmitro44/relper/internal/repositories/users"
	"github.com/google/uuid"
)

type UserService struct {
	repo      usersrepo.Repository
	rolesRepo rolesrepo.Repository
}

func NewUserService(repo usersrepo.Repository, rolesRepo rolesrepo.Repository) *UserService {
	return &UserService{
		repo:      repo,
		rolesRepo: rolesRepo,
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
	userEntity, err := u.repo.GetByEmail(ctx, userentity.NormalizeEmail(email))
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
		Email:        userentity.NormalizeEmail(input.Email),
		PasswordHash: input.PasswordHash,
		Nickname:     input.Nickname,
		Habits:       input.Habits,
		Contacts:     input.Contacts,
	})
	if err != nil {
		return nil, apperr.NewInternalError("failed to create user", err)
	}
	err = u.rolesRepo.Add(ctx, userEntity.ID, []string{"user"})
	if err != nil {
		return nil, apperr.NewInternalError("failed to add role to user", err)
	}
	slog.InfoContext(ctx, "user created", "user_id", userEntity.ID)

	return userEntity, nil
}

func (u *UserService) Update(ctx context.Context, id, callerID uuid.UUID, input UpdateInput) (*userentity.User, error) {
	if id != callerID {
		return nil, apperr.NewForbiddenError("cannot modify another user")
	}
	if input.Email != nil {
		normEmail := userentity.NormalizeEmail(*input.Email)
		input.Email = &normEmail
	}

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

func (u *UserService) Delete(ctx context.Context, id, callerID uuid.UUID) error {
	if id != callerID {
		return apperr.NewForbiddenError("cannot delete another user")
	}

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

func (u *UserService) GetRoles(ctx context.Context, id uuid.UUID) ([]string, error) {
	roles, err := u.rolesRepo.GetByUserID(ctx, id)
	if err != nil {
		return nil, apperr.NewInternalError("failed to get user roles", err)
	}

	return roles, nil
}

func (u *UserService) SetRoles(ctx context.Context, id uuid.UUID, codes []string) ([]string, error) {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, usersrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("user not found")
		}
		return nil, apperr.NewInternalError("failed to get user", err)
	}

	err := u.rolesRepo.ReplaceForUser(ctx, id, codes)
	if err != nil {
		if errors.Is(err, rolesrepo.ErrUnknownRole) {
			return nil, apperr.NewBadRequestError("unknown role code")
		}
		return nil, apperr.NewInternalError("failed to set user roles", err)
	}
	slog.InfoContext(ctx, "user roles set", "user_id", id, "roles", codes)

	return u.GetRoles(ctx, id)
}
