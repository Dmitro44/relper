package dto

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       uuid.UUID
	Email    string
	Nickname string
	Habits   []string
	Contacts ContactsInput
}

type ContactsInput struct {
	Telegram string `json:"telegram,omitempty" minLength:"1" doc:"Telegram handle"`
	Email    string `json:"email,omitempty" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"Contact email"`
}
