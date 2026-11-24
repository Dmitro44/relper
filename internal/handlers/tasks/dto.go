package tasks

import (
	"time"

	"github.com/Dmitro44/relper/internal/entities/task"
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Body struct {
		Title       string    `json:"title" minLength:"1"`
		Description string    `json:"description,omitempty"`
		PlaceID     uuid.UUID `json:"place_id"`
		Done        bool      `json:"done,omitempty"`
	}
}

type UpdateTaskRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Title       *string    `json:"title,omitempty" minLength:"1"`
		Description *string    `json:"description,omitempty"`
		PlaceID     *uuid.UUID `json:"place_id,omitempty"`
		Done        *bool      `json:"done,omitempty"`
	}
}

type GetTaskRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteTaskRequest struct {
	ID uuid.UUID `path:"id"`
}

type TaskItem struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PlaceID     uuid.UUID `json:"place_id"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskResponse struct {
	Body TaskItem
}

func toTaskItem(t task.Task) TaskItem {
	return TaskItem{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		PlaceID:     t.PlaceID,
		Done:        t.Done,
		CreatedAt:   t.CreatedAt,
	}
}
