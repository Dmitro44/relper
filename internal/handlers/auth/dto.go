package auth

import (
	"net/http"

	"github.com/Dmitro44/relper/internal/handlers/dto"
)

type RegisterRequest struct {
	Body struct {
		Email    string             `json:"email" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"User email"`
		Password string             `json:"password" doc:"User password"`
		Nickname string             `json:"nickname" minLength:"1" maxLength:"50" doc:"Display name"`
		Habits   []string           `json:"habits,omitempty" doc:"User habits"`
		Contacts *dto.ContactsInput `json:"contacts" doc:"Contact info"`
	}
}

type LoginRequest struct {
	Body struct {
		Email    string `json:"email" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"User email"`
		Password string `json:"password" doc:"User password"`
	}
}

type RefreshRequest struct {
	RefreshToken string `cookie:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `cookie:"refresh_token"`
}

type LogoutResponse struct {
	ClearCookie http.Cookie `header:"Set-Cookie"`
}

type MeResponse struct {
	Body *dto.UserResponse
}

type AuthReponseBody struct {
	User        *dto.UserResponse `json:"user"`
	AccessToken string            `json:"access_token"`
}
type AuthResponse struct {
	Body          AuthReponseBody
	RefreshCookie http.Cookie `header:"Set-Cookie"`
}
