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
	UserID      uuid.UUID
	Title       string
	Description string
	PlaceID     uuid.UUID
	Done        bool
}

type UpdateInput struct {
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

func (t *TaskService) GetByID(ctx context.Context, id, callerID uuid.UUID, admin bool) (*task.Task, error) {
	taskEntity, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tasksrepo.ErrNotFound) {
			return nil, apperr.NewNotFoundError("task not found")
		}
		return nil, apperr.NewInternalError("failed to get task", err)
	}

	if taskEntity.UserID != callerID && !admin {
		slog.InfoContext(ctx, "task access denied", "task_id", id, "caller_id", callerID)
		return nil, apperr.NewNotFoundError("task not found")
	}

	return taskEntity, nil
}

func (t *TaskService) Create(ctx context.Context, input CreateInput) (*task.Task, error) {
	taskEntity, err := t.repo.Create(ctx, tasksrepo.CreateInput{
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

func (t *TaskService) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input UpdateInput) (*task.Task, error) {
	taskEntity, err := t.repo.Update(ctx, id, userID, admin, tasksrepo.UpdateInput{
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

func (t *TaskService) Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error {
	err := t.repo.Delete(ctx, id, userID, admin)
	if err != nil {
		if errors.Is(err, tasksrepo.ErrNotFound) {
			return apperr.NewNotFoundError("task not found")
		}
		return apperr.NewInternalError("failed to delete task", err)
	}
	slog.InfoContext(ctx, "task deleted", "task_id", id)

	return nil
}
