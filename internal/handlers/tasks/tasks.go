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

func (h *TaskHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Tasks"}
	})

	huma.Get(grp, "/tasks/{id}", h.HandleGetByID)
	huma.Patch(grp, "/tasks/{id}", h.HandleUpdate)
	huma.Delete(grp, "/tasks/{id}", h.HandleDelete)

	huma.Register(grp, huma.Operation{
		OperationID:   "create-task",
		Method:        http.MethodPost,
		Path:          "/tasks",
		Summary:       "Create a task",
		DefaultStatus: 201,
	}, h.HandleCreate)
}

func (h *TaskHandler) HandleGetByID(ctx context.Context, in *GetTaskRequest) (*TaskResponse, error) {
	t, err := h.svc.GetByID(ctx, in.ID)
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
	t, err := h.svc.Update(ctx, in.ID, userID, tasksvc.UpdateInput{
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
	err := h.svc.Delete(ctx, in.ID, userID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
