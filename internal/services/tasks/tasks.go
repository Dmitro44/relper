package tasks

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/task"
	tasksrepo "github.com/Dmitro44/relper/internal/repositories/tasks"
	"github.com/google/uuid"
)

type TaskService struct {
	repo tasksrepo.Repository
}

type CreateInput struct {
	LocationID  uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	PlaceID     uuid.UUID
	Done        bool
}

type UpdateInput struct {
	LocationID  *uuid.UUID
	UserID      *uuid.UUID
	Title       *string
	Description *string
	PlaceID     *uuid.UUID
	Done        *bool
}

func (t *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	taskEntity, err := t.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return taskEntity, nil
}

func (t *TaskService) Create(ctx context.Context, input CreateInput) (*task.Task, error) {
	taskEntity, err := t.repo.Create(ctx, tasksrepo.CreateInput{
		LocationID:  input.LocationID,
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		PlaceID:     input.PlaceID,
		Done:        input.Done,
	})
	if err != nil {
		return nil, err
	}

	return taskEntity, nil
}

func (t *TaskService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*task.Task, error) {
	taskEntity, err := t.repo.Update(ctx, id, tasksrepo.UpdateInput{
		LocationID:  input.LocationID,
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		PlaceID:     input.PlaceID,
		Done:        input.Done,
	})
	if err != nil {
		return nil, err
	}

	return taskEntity, nil
}
