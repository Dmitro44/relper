package location

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/google/uuid"
)

type Location struct {
	ID          uuid.UUID
	Country     string
	City        string
	Coordinates coordinates.Coordinates
	CreatedAt   time.Time
}
