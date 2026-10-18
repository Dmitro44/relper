package roles

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, id uuid.UUID, codes []string) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
	ReplaceForUser(ctx context.Context, userID uuid.UUID, codes []string) error
}
