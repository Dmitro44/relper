package model

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Nickname  string
	Habits    []string
	Contacts  user.Contacts
	CreatedAt time.Time
}
