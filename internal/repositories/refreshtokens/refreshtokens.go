package refreshtokens

import (
	"context"
	"time"

	"github.com/Dmitro44/relper/internal/entities/refreshtoken"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	GetByHash(ctx context.Context, tokenHash string) (*refreshtoken.RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
}

type CreateInput struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

type UpdateInput struct {
	UserID     *uuid.UUID
	TokenHash  *string
	ExpiresAt  *time.Time
	RevokedAt  *time.Time
	ReplacedBy *uuid.UUID
}
