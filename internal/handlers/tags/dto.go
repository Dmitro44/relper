package tags

import (
	"time"

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

type TagItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type TagResponse struct {
	Body TagItem
}

func toTagItem(t tag.Tag) TagItem {
	return TagItem{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
	}
}
