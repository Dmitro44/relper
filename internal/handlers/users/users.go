package users

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/user"
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

	// Move in mapper
	dto := &dto.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,
		Habits:   u.Habits,
		Contacts: dto.ContactsInput(u.Contacts),
	}

	return &UserResponse{Body: dto}, nil
}

func (h *UserHandler) HandleUpdate(ctx context.Context, in *UpdateUserRequest) (*UserResponse, error) {
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

	u, err := h.svc.Update(ctx, in.ID, updateInput)
	if err != nil {
		return nil, httperr.Map(err)
	}

	// Move in mapper
	dto := &dto.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,
		Habits:   u.Habits,
		Contacts: dto.ContactsInput(u.Contacts),
	}

	return &UserResponse{Body: dto}, nil
}

func (h *UserHandler) HandleDelete(ctx context.Context, in *DeleteUserRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
