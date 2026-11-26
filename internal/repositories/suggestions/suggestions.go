package suggestions

import (
	"context"
	"fmt"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/suggestion"
	"github.com/google/uuid"
)

var (
	ErrNotFound   = fmt.Errorf("suggestion not found")
	ErrNotPending = fmt.Errorf("suggestion is not pending")
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) (*suggestion.Suggestion, error)
	ListMine(ctx context.Context, userID uuid.UUID) ([]suggestion.Suggestion, error)
	ListByStatus(ctx context.Context, status string) ([]suggestion.Suggestion, error)
	Approve(ctx context.Context, id, reviewerID uuid.UUID) (*suggestion.Suggestion, error)
	Reject(ctx context.Context, id, reviewerID uuid.UUID, reason string) (*suggestion.Suggestion, error)
}

type CreateInput struct {
	UserID      uuid.UUID
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Description string
	Coordinates coordinates.Coordinates
}
