package users

import (
	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Body struct {
		Email    string        `json:"email" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"User email"`
		Nickname string        `json:"nickname" minLength:"1" maxLength:"50" doc:"Display name"`
		Habits   []string      `json:"habits,omitempty" doc:"User habits"`
		Contacts ContactsInput `json:"contacts" doc:"Contact info"`
	}
}

type ContactsInput struct {
	Telegram string `json:"telegram,omitempty" minLength:"1" doc:"Telegram handle"`
	Email    string `json:"email,omitempty" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"Contact email"`
}

type UpdateUserRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Email    *string        `json:"email,omitempty" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"User email"`
		Nickname *string        `json:"nickname,omitempty" minLength:"1" maxLength:"50" doc:"Display name"`
		Habits   *[]string      `json:"habits,omitempty" doc:"User habits"`
		Contacts *ContactsInput `json:"contacts,omitempty" doc:"Contact info"`
	}
}

type GetUserRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `path:"id"`
}

type UserResponse struct {
	Body *user.User
}
