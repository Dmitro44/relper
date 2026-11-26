package audit

import (
	"time"

	"github.com/google/uuid"
)

type ListAuditRequest struct {
	Source     string    `query:"source"`
	Action     string    `query:"action"`
	EntityType string    `query:"entity_type"`
	ActorID    uuid.UUID `query:"actor_id"`
	Limit      int       `query:"limit" default:"500" minimum:"1" maximum:"500"`
}

type AuditItem struct {
	ID            uuid.UUID  `json:"id"`
	Source        string     `json:"source"`
	Action        string     `json:"action"`
	EntityType    string     `json:"entity_type"`
	EntityID      *uuid.UUID `json:"entity_id,omitempty"`
	ActorID       *uuid.UUID `json:"actor_id,omitempty"`
	ActorNickname string     `json:"actor_nickname,omitempty"`
	Details       any        `json:"details"`
	CreatedAt     time.Time  `json:"created_at"`
}

type AuditResponse struct {
	Body []AuditItem
}
