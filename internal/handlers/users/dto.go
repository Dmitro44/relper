package users

import (
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Email    *string            `json:"email,omitempty" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"User email"`
		Nickname *string            `json:"nickname,omitempty" minLength:"1" maxLength:"50" doc:"Display name"`
		Habits   *[]string          `json:"habits,omitempty" doc:"User habits"`
		Contacts *dto.ContactsInput `json:"contacts,omitempty" doc:"Contact info"`
	}
}

type GetUserRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `path:"id"`
}

type UserResponse struct {
	Body *dto.UserResponse
}
