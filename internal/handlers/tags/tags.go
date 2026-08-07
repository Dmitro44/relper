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

func (h *TagHandler) Register(api huma.API) {
	grp := huma.NewGroup(api)
	grp.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Tags"}
	})

	huma.Get(grp, "/tags/{id}", h.HandleGetByID)
	huma.Patch(grp, "/tags/{id}", h.HandleUpdate)
	huma.Delete(grp, "/tags/{id}", h.HandleDelete)

	huma.Register(grp, huma.Operation{
		OperationID:   "create-tag",
		Method:        http.MethodPost,
		Path:          "/tags",
		Summary:       "Create a tag",
		DefaultStatus: 201,
	}, h.HandleCreate)
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
