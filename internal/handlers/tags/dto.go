package tags

import (
	"github.com/Dmitro44/relper/internal/entities/tag"
	"github.com/google/uuid"
)

type CreateTagRequest struct {
	Body struct {
		Name string `json:"name" minLength:"1"`
		Slug string `json:"slug" pattern:"^[a-z0-9-]+$"`
	}
}

type UpdateTagRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name *string `json:"name,omitempty" minLength:"1"`
		Slug *string `json:"slug,omitempty" pattern:"^[a-z0-9-]+$"`
	}
}

type GetTagRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteTagRequest struct {
	ID uuid.UUID `path:"id"`
}

type TagResponse struct {
	Body *tag.Tag
}
