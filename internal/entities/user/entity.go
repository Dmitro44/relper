package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Email    string
	Nickname string
	Habits   []string
	Contacts Contacts
}

type Contacts struct {
	Telegram string `json:"telegram,omitempty"`
	Email    string `json:"email,omitempty"`
}
