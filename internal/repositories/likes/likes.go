package likes

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	LikeTip(ctx context.Context, userID, tipID uuid.UUID) error
	LikeComment(ctx context.Context, userID, commentID uuid.UUID) error

	UnlikeTip(ctx context.Context, userID, tipID uuid.UUID) error
	UnlikeComment(ctx context.Context, userID, commentID uuid.UUID) error

	LikedTip(ctx context.Context, userID, tipID uuid.UUID) (bool, error)
	LikedComment(ctx context.Context, userID, commentID uuid.UUID) (bool, error)

	CountTip(ctx context.Context, tipID uuid.UUID) (int64, error)
	CountComment(ctx context.Context, commentID uuid.UUID) (int64, error)
}
