package auth

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
	"github.com/Dmitro44/relper/internal/entities/user"
	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	"github.com/Dmitro44/relper/internal/handlers/users"
	authsvc "github.com/Dmitro44/relper/internal/services/auth"
	usersvc "github.com/Dmitro44/relper/internal/services/users"
	"github.com/danielgtaylor/huma/v2"
)

type AuthHandler struct {
	svc     *authsvc.AuthService
	userSvc *usersvc.UserService
}

func NewAuthHandler(svc *authsvc.AuthService, userSvc *usersvc.UserService) *AuthHandler {
	return &AuthHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

func (h *AuthHandler) Register(public huma.API, protected huma.API) {
	huma.Post(public, "/auth/login", h.HandleLogin)
	huma.Post(public, "/auth/refresh", h.HandleRefresh)

	huma.Register(public, huma.Operation{
		OperationID:   "register-user",
		Method:        http.MethodPost,
		Path:          "/auth/register",
		Summary:       "Register a user",
		DefaultStatus: http.StatusCreated,
	}, h.HandleRegister)

	huma.Register(public, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/auth/logout",
		Summary:       "Logout",
		DefaultStatus: http.StatusNoContent,
	}, h.HandleLogout)

	huma.Get(protected, "/auth/me", h.HandleMe)
}

func (h *AuthHandler) HandleRegister(ctx context.Context, in *RegisterRequest) (*AuthResponse, error) {
	input := authsvc.CreateInput{
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

	return &AuthResponse{
		Body: AuthReponseBody{
			User:        users.MapUserToDto(user),
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

	return &AuthResponse{
		Body: AuthReponseBody{
			User:        users.MapUserToDto(user),
			AccessToken: pair.AccessToken,
		},
		RefreshCookie: refreshCookie(pair.RefreshToken),
	}, nil
}

func (h *AuthHandler) HandleLogout(ctx context.Context, in *LogoutRequest) (*LogoutResponse, error) {
	if err := h.svc.Logout(ctx, in.RefreshToken); err != nil {
		return nil, httperr.Map(err)
	}

	return &LogoutResponse{ClearCookie: clearRefreshCookie()}, nil
}

func (h *AuthHandler) HandleMe(ctx context.Context, _ *struct{}) (*MeResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(apperr.NewUnauthorizedError("unauthenticated"))
	}

	u, err := h.userSvc.GetByID(ctx, userID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &MeResponse{Body: users.MapUserToDto(u)}, nil
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
		Path:     "/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   15 * 24 * 60 * 60,
	}
}

func clearRefreshCookie() http.Cookie {
	c := refreshCookie("")
	c.MaxAge = -1
	return c
}
