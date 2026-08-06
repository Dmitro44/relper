package tasks

import (
	"github.com/Dmitro44/relper/internal/entities/task"
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Body struct {
		LocationID  uuid.UUID `json:"location_id"`
		UserID      uuid.UUID `json:"user_id"`
		Title       string    `json:"title" minLength:"1"`
		Description string    `json:"description,omitempty"`
		PlaceID     uuid.UUID `json:"place_id"`
		Done        bool      `json:"done,omitempty"`
	}
}

type UpdateTaskRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		LocationID  *uuid.UUID `json:"location_id,omitempty"`
		UserID      *uuid.UUID `json:"user_id,omitempty"`
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

type TaskResponse struct {
	Body *task.Task
}
