package tags

import (
	"context"
	"net/http"

	"github.com/Dmitro44/relper/internal/handlers/httperr"
	tagsvc "github.com/Dmitro44/relper/internal/services/tags"
	"github.com/danielgtaylor/huma/v2"
)

type TagHandler struct {
	svc *tagsvc.TagService
}

func NewTagHandler(svc *tagsvc.TagService) *TagHandler {
	return &TagHandler{
		svc: svc,
	}
}

func (h *TagHandler) Register(public, admin huma.API) {
	huma.Register(public, huma.Operation{
		OperationID:   "get-tag",
		Method:        http.MethodGet,
		Path:          "/tags/{id}",
		Summary:       "Get a tag",
		Tags:          []string{"Tags"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetByID)

	huma.Register(admin, huma.Operation{
		OperationID:   "create-tag",
		Method:        http.MethodPost,
		Path:          "/tags",
		Summary:       "Create a tag",
		Tags:          []string{"Tags"},
		DefaultStatus: http.StatusCreated,
	}, h.HandleCreate)

	huma.Register(admin, huma.Operation{
		OperationID:   "update-tag",
		Method:        http.MethodPatch,
		Path:          "/tags/{id}",
		Summary:       "Update a tag",
		Tags:          []string{"Tags"},
		DefaultStatus: http.StatusOK,
	}, h.HandleUpdate)

	huma.Register(admin, huma.Operation{
		OperationID:   "delete-tag",
		Method:        http.MethodDelete,
		Path:          "/tags/{id}",
		Summary:       "Delete a tag",
		Tags:          []string{"Tags"},
		DefaultStatus: http.StatusNoContent,
	}, h.HandleDelete)
}

func (h *TagHandler) HandleGetByID(ctx context.Context, in *GetTagRequest) (*TagResponse, error) {
	t, err := h.svc.GetByID(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TagResponse{Body: t}, nil
}

func (h *TagHandler) HandleCreate(ctx context.Context, in *CreateTagRequest) (*TagResponse, error) {
	t, err := h.svc.Create(ctx, tagsvc.CreateInput{
		Name: in.Body.Name,
		Slug: in.Body.Slug,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TagResponse{Body: t}, nil
}

func (h *TagHandler) HandleUpdate(ctx context.Context, in *UpdateTagRequest) (*TagResponse, error) {
	t, err := h.svc.Update(ctx, in.ID, tagsvc.UpdateInput{
		Name: in.Body.Name,
		Slug: in.Body.Slug,
	})
	if err != nil {
		return nil, httperr.Map(err)
	}

	return &TagResponse{Body: t}, nil
}

func (h *TagHandler) HandleDelete(ctx context.Context, in *DeleteTagRequest) (*struct{}, error) {
	err := h.svc.Delete(ctx, in.ID)
	if err != nil {
		return nil, httperr.Map(err)
	}

	return nil, nil
}
