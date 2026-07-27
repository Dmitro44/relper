package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Nickname  string
	Habits    []string
	Contacts  Contacts
	CreatedAt time.Time
}

type Contacts struct {
	Telegram string `json:"telegram,omitempty"`
	Email    string `json:"email,omitempty"`
}
