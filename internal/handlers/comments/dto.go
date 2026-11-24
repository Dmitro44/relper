package comments

import (
	"time"

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

type ListCommentsRequest struct {
	ID uuid.UUID `path:"id"`
}

type CommentItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentsResponse struct {
	Body []CommentItem
}

func toCommentItems(items []comment.WithAuthor) []CommentItem {
	out := make([]CommentItem, 0, len(items))
	for _, item := range items {
		out = append(out, CommentItem{
			ID:        item.ID,
			UserID:    item.UserID,
			Nickname:  item.Nickname,
			Content:   item.Content,
			CreatedAt: item.CreatedAt,
		})
	}

	return out
}
