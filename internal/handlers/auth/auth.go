package auth

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/Dmitro44/relper/internal/services/auth"
	"github.com/danielgtaylor/huma/v2"
)

type AuthHandler struct {
	svc *auth.AuthService
}

func NewAuthHandler(svc *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Register(api huma.API) {
	huma.Post(api, "/auth/login", h.HandleLogin)

	huma.Register(api, huma.Operation{
		OperationID:   "register-user",
		Method:        http.MethodPost,
		Path:          "/auth/register",
		Summary:       "Register a user",
		DefaultStatus: http.StatusCreated,
	}, h.HandleRegister)
}

func (h *AuthHandler) HandleRegister(ctx context.Context, in *RegisterRequest) (*AuthResponse, error) {
	input := auth.CreateInput{
		Email:    in.Body.Email,
		Password: in.Body.Password,
		Nickname: in.Body.Nickname,
		Habits:   in.Body.Habits,
	}

	if in.Body.Contacts != nil {
		input.Contacts = user.Contacts{
			Telegram: in.Body.Contacts.Telegram,
			Email:    in.Body.Contacts.Email,
		}
	}

	user, token, err := h.svc.Register(ctx, input)
	if err != nil {
		return nil, err
	}

	// Move in mapper
	dto := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Habits:   user.Habits,
		Contacts: dto.ContactsInput(user.Contacts),
	}

	return &AuthResponse{
		User:        dto,
		AccessToken: token,
	}, nil
}

func (h *AuthHandler) HandleLogin(ctx context.Context, in *LoginRequest) (*AuthResponse, error) {
	user, token, err := h.svc.Login(ctx, in.Body.Email, in.Body.Password)
	if err != nil {
		return nil, err
	}

	dto := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Habits:   user.Habits,
		Contacts: dto.ContactsInput(user.Contacts),
	}

	return &AuthResponse{
		User:        dto,
		AccessToken: token,
	}, nil
}
