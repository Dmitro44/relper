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
		TargetType string    `json:"target_type" minLength:"1" doc:"Type of the target (e.g. tip, place)"`
		TargetID   uuid.UUID `json:"target_id"`
		UserID     uuid.UUID `json:"user_id"`
		Content    string    `json:"content" minLength:"1"`
	}
}

type UpdateCommentRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		TargetType *string    `json:"target_type,omitempty"`
		TargetID   *uuid.UUID `json:"target_id,omitempty"`
		UserID     *uuid.UUID `json:"user_id,omitempty"`
		Content    *string    `json:"content,omitempty" minLength:"1"`
	}
}

type CommentResponse struct {
	Body *comment.Comment
}
