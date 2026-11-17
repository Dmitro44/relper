package likes

import (
	"github.com/Dmitro44/relper/internal/entities/like"
	"github.com/google/uuid"
)

type LikeTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type UnlikeTipRequest struct {
	ID uuid.UUID `path:"id"`
}

type GetTipLikesRequest struct {
	ID uuid.UUID `path:"id"`
}

type LikeCommentRequest struct {
	ID uuid.UUID `path:"id"`
}

type UnlikeCommentRequest struct {
	ID uuid.UUID `path:"id"`
}

type GetCommentLikesRequest struct {
	ID uuid.UUID `path:"id"`
}

type SummaryResponse struct {
	Body *like.LikeSummary
}
