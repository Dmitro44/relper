package tasks

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/auth"
	"github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/handlers/httperr"
	tasksvc "github.com/Dmitro44/relper/internal/services/tasks"
	"github.com/danielgtaylor/huma/v2"
)

type TaskHandler struct {
	svc *tasksvc.TaskService
}

func NewTaskHandler(svc *tasksvc.TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

func (h *TaskHandler) Register(public huma.API) {
	huma.Register(public, huma.Operation{
		OperationID:   "get-task",
		Method:        http.MethodGet,
		Path:          "/tasks/{id}",
		Summary:       "Get a task",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetByID)

	huma.Register(public, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost,
		Path:          "/tasks",
		Summary:       "Create a task",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)

	huma.Register(public, huma.Operation{
		OperationID:   "update-task",
		Method:        http.MethodPatch,
		Path:          "/tasks/{id}",
		Summary:       "Update a task",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUpdate)

	huma.Register(public, huma.Operation{
		OperationID:   "delete-task",
		Method:        http.MethodDelete,
		Path:          "/tasks/{id}",
		Summary:       "Delete a task",
		Tags:          []string{"Tasks"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)
}

func (h *TaskHandler) HandleGetByID(ctx context.Context, in *GetTaskRequest) (*TaskResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	t, err := h.svc.GetByID(ctx, in.ID, userID, auth.HasRole(ctx, "admin"))
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TaskResponse{Body: t}, nil
}

func (h *TaskHandler) HandleCreate(ctx context.Context, in *CreateTaskRequest) (*TaskResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	t, err := h.svc.Create(ctx, tasksvc.CreateInput{
		LocationID:  in.Body.LocationID,
		UserID:      userID,
		Title:       in.Body.Title,
		Description: in.Body.Description,
		PlaceID:     in.Body.PlaceID,
		Done:        in.Body.Done,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TaskResponse{Body: t}, nil
}

func (h *TaskHandler) HandleUpdate(ctx context.Context, in *UpdateTaskRequest) (*TaskResponse, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	t, err := h.svc.Update(ctx, in.ID, userID, auth.HasRole(ctx, "admin"), tasksvc.UpdateInput{
		LocationID:  in.Body.LocationID,
		Title:       in.Body.Title,
		Description: in.Body.Description,
		PlaceID:     in.Body.PlaceID,
		Done:        in.Body.Done,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TaskResponse{Body: t}, nil
}

func (h *TaskHandler) HandleDelete(ctx context.Context, in *DeleteTaskRequest) (*struct{}, error) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return nil, httperr.Map(errors.NewUnauthorizedError("invalid user id from context"))
	}
	err := h.svc.Delete(ctx, in.ID, userID, auth.HasRole(ctx, "admin"))
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
