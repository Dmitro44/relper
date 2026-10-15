package users

import (
	"context"

	"github.com/Dmitro44/relper/internal/authn"
	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	usersvc "github.com/Dmitro44/relper/internal/services/users"
	"github.com/danielgtaylor/huma/v2"
)

type UserHandler struct {
	svc *usersvc.UserService
}

func NewUserHandler(svc *usersvc.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func MapUserToDto(u *user.User) *dto.UserResponse {
	if u == nil {
		return nil
	}

	return &dto.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,
		Habits:   u.Habits,
		Contacts: dto.ContactsInput(u.Contacts),
	}
}

func (h *UserHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Users"}
	})

	huma.Get(grp, "/users/{id}", h.HandleGetByID)
	huma.Patch(grp, "/users/{id}", h.HandleUpdate)
	huma.Delete(grp, "/users/{id}", h.HandleDelete)
}

func (h *UserHandler) HandleGetByID(ctx context.Context, in *GetUserRequest) (*UserResponse, error) {
	u, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &UserResponse{Body: MapUserToDto(u)}, nil
}

func (h *UserHandler) HandleUpdate(ctx context.Context, in *UpdateUserRequest) (*UserResponse, error) {
	callerID, ok := authn.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	var contacts *user.Contacts
	if in.Body.Contacts != nil {
		contacts = &user.Contacts{
			Telegram: in.Body.Contacts.Telegram,
			Email:    in.Body.Contacts.Email,
		}
	}
	updateInput := usersvc.UpdateInput{
		Email:    in.Body.Email,
		Nickname: in.Body.Nickname,
		Habits:   in.Body.Habits,
		Contacts: contacts,
	}

	u, err := h.svc.Update(ctx, in.ID, callerID, updateInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &UserResponse{Body: MapUserToDto(u)}, nil
}

func (h *UserHandler) HandleDelete(ctx context.Context, in *DeleteUserRequest) (*struct{}, error) {
	callerID, ok := authn.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}

	err := h.svc.Delete(ctx, in.ID, callerID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
