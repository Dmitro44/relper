package auth

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/entities/user"
	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
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
	huma.Post(api, "/auth/refresh", h.HandleRefresh)

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

	user, pair, err := h.svc.Register(ctx, input)
	if err != nil {
		return nil, httperr.Map(err)
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
		Body: AuthReponseBody{
			User:        dto,
			AccessToken: pair.AccessToken,
		},
		RefreshCookie: refreshCookie(pair.RefreshToken),
	}, nil
}

func (h *AuthHandler) HandleLogin(ctx context.Context, in *LoginRequest) (*AuthResponse, error) {
	user, pair, err := h.svc.Login(ctx, in.Body.Email, in.Body.Password)
	if err != nil {
		return nil, httperr.Map(err)
	}

	dto := &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Habits:   user.Habits,
		Contacts: dto.ContactsInput(user.Contacts),
	}

	return &AuthResponse{
		Body: AuthReponseBody{
			User:        dto,
			AccessToken: pair.AccessToken,
		},
		RefreshCookie: refreshCookie(pair.RefreshToken),
	}, nil
}

func (h *AuthHandler) HandleRefresh(ctx context.Context, in *RefreshRequest) (*AuthResponse, error) {
	if in.RefreshToken == "" {
		return nil, httperr.Map(apperr.NewUnauthorizedError("refresh token required"))
	}

	pair, err := h.svc.Refresh(ctx, in.RefreshToken)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &AuthResponse{
		Body:          AuthReponseBody{AccessToken: pair.AccessToken},
		RefreshCookie: refreshCookie(pair.RefreshToken),
	}, nil
}

func refreshCookie(token string) http.Cookie {
	return http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   15 * 24 * 60 * 60,
	}
}
