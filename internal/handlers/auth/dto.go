package auth

import (
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

type AuthResponse struct {
	User        *dto.UserResponse
	AccessToken string
}
