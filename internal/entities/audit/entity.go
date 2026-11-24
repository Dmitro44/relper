package audit

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID
	Source        string
	Action        string
	EntityType    string
	EntityID      *uuid.UUID
	ActorID       *uuid.UUID
	ActorNickname string
	Details       []byte
	CreatedAt     time.Time
}

type EventInput struct {
	Action     string
	EntityType string
	EntityID   *uuid.UUID
	ActorID    *uuid.UUID
	Details    any
}

type Filter struct {
	Source     string
	Action     string
	EntityType string
	ActorID    *uuid.UUID
	Limit      int
}
