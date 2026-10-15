package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Nickname     string
	Habits       []string
	Contacts     Contacts
	CreatedAt    time.Time
}

type Contacts struct {
	Telegram string `json:"telegram,omitempty"`
	Email    string `json:"email,omitempty"`
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
