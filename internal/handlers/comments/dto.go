package comments

import (
	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/google/uuid"
)

type GetCommentRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteCommentRequest struct {
	ID uuid.UUID `path:"id"`
}

type CreateCommentRequest struct {
	Body struct {
		PlaceID *uuid.UUID `json:"place_id,omitempty"`
		TipID   *uuid.UUID `json:"tip_id,omitempty"`
		Content string     `json:"content" minLength:"1"`
	}
}

type UpdateCommentRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Content *string `json:"content,omitempty" minLength:"1"`
	}
}

type CommentResponse struct {
	Body *comment.Comment
}
