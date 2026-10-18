package dto

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       uuid.UUID      `json:"id"`
	Email    string         `json:"email"`
	Nickname string         `json:"nickname"`
	Habits   []string       `json:"habits,omitempty"`
	Contacts ContactsOutput `json:"contacts"`
}

type ContactsOutput struct {
	Telegram string `json:"telegram,omitempty" doc:"Telegram handle"`
	Email    string `json:"email,omitempty" doc:"Contact email"`
}

type ContactsInput struct {
	Telegram string `json:"telegram,omitempty" minLength:"1" doc:"Telegram handle"`
	Email    string `json:"email,omitempty" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"Contact email"`
}
