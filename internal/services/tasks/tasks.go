package tasks

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Dmitro44/relper/internal/entities/task"
	apperr "github.com/Dmitro44/relper/internal/errors"
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
	Title       *string
	Description *string
	PlaceID     *uuid.UUID
	Done        *bool
}

func NewTaskService(repo tasksrepo.Repository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (t *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	taskEntity, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tasksrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("task not found")
		}
		return nil, apperr.NewInternalError("failed to get task", err)
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
		return nil, apperr.NewInternalError("failed to create task", err)
	}
	slog.InfoContext(ctx, "task created", "task_id", taskEntity.ID)

	return taskEntity, nil
}

func (t *TaskService) Update(ctx context.Context, id, userID uuid.UUID, input UpdateInput) (*task.Task, error) {
	taskEntity, err := t.repo.Update(ctx, id, userID, tasksrepo.UpdateInput{
		LocationID:  input.LocationID,
		Title:       input.Title,
		Description: input.Description,
		PlaceID:     input.PlaceID,
		Done:        input.Done,
	})
	if err != nil {
		if errors.Is(err, tasksrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("task not found")
		}
		return nil, apperr.NewInternalError("failed to update task", err)
	}
	slog.InfoContext(ctx, "task updated", "task_id", taskEntity.ID)

	return taskEntity, nil
}

func (t *TaskService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	err := t.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, tasksrepo.ErrNotFound) {
			return apperr.NewNotFoundError("task not found")
		}
		return apperr.NewInternalError("failed to delete task", err)
	}
	slog.InfoContext(ctx, "task deleted", "task_id", id)

	return nil
}
