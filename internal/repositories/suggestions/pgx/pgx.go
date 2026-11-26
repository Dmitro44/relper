package pgx

import (
	"context"
	"errors"
	"time"

	"github.com/Dmitro44/relper/internal/entities/suggestion"
	"github.com/Dmitro44/relper/internal/repositories/suggestions"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ suggestions.Repository = (*Pgx)(nil)

const suggestionColumns = `s.id, s.user_id, u.nickname, s.location_id, s.name, s.category,
	s.address, s.phone, s.website_url, s.description, s.latitude, s.longitude,
	s.status, s.rejection_reason, s.reviewed_by, s.reviewed_at, s.place_id, s.created_at`

const suggestionFrom = `from place_suggestions s join users u on u.id = s.user_id`

type Pgx struct {
	pool *pgxpool.Pool
}

func NewSuggestionsRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{pool: pool}
}

func (p *Pgx) Create(ctx context.Context, input suggestions.CreateInput) (*suggestion.Suggestion, error) {
	query := `
insert into place_suggestions(user_id, location_id, name, category, address, phone, website_url, description, latitude, longitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning id;
`

	var id uuid.UUID
	err := p.pool.QueryRow(
		ctx, query,
		input.UserID,
		input.LocationID,
		input.Name,
		input.Category,
		input.Address,
		input.Phone,
		input.WebsiteURL,
		input.Description,
		input.Coordinates.Latitude,
		input.Coordinates.Longitude,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return p.getByID(ctx, id)
}

func (p *Pgx) getByID(ctx context.Context, id uuid.UUID) (*suggestion.Suggestion, error) {
	query := `select ` + suggestionColumns + ` ` + suggestionFrom + ` where s.id = $1`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[suggestion.Suggestion])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, suggestions.ErrNotFound
		}
		return nil, err
	}

	return &item, nil
}

func (p *Pgx) ListMine(ctx context.Context, userID uuid.UUID) ([]suggestion.Suggestion, error) {
	query := `select ` + suggestionColumns + ` ` + suggestionFrom + ` where s.user_id = $1 order by s.created_at desc`

	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[suggestion.Suggestion])
}

func (p *Pgx) ListByStatus(ctx context.Context, status string) ([]suggestion.Suggestion, error) {
	query := `select ` + suggestionColumns + ` ` + suggestionFrom + ` where s.status = $1 order by s.created_at desc`

	rows, err := p.pool.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[suggestion.Suggestion])
}

func (p *Pgx) Approve(ctx context.Context, id, reviewerID uuid.UUID) (*suggestion.Suggestion, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	item, err := pendingForUpdate(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	query := `
insert into places(location_id, name, category, address, phone, website_url, latitude, longitude, description)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning id;
`

	var placeID uuid.UUID
	err = tx.QueryRow(
		ctx, query,
		item.LocationID,
		item.Name,
		item.Category,
		item.Address,
		item.Phone,
		item.WebsiteURL,
		item.Coordinates.Latitude,
		item.Coordinates.Longitude,
		item.Description,
	).Scan(&placeID)
	if err != nil {
		return nil, err
	}

	var reviewedAt time.Time
	query = `
update place_suggestions
set status = $2, reviewed_by = $3, reviewed_at = now(), place_id = $4, rejection_reason = null
where id = $1
returning reviewed_at;
`
	err = tx.QueryRow(ctx, query, id, suggestion.StatusApproved, reviewerID, placeID).Scan(&reviewedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	item.Status = suggestion.StatusApproved
	item.ReviewedBy = &reviewerID
	item.ReviewedAt = &reviewedAt
	item.PlaceID = &placeID
	item.RejectionReason = nil

	return item, nil
}

func (p *Pgx) Reject(ctx context.Context, id, reviewerID uuid.UUID, reason string) (*suggestion.Suggestion, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	item, err := pendingForUpdate(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	var reviewedAt time.Time
	query := `
update place_suggestions
set status = $2, reviewed_by = $3, reviewed_at = now(), rejection_reason = $4
where id = $1
returning reviewed_at;
`
	err = tx.QueryRow(ctx, query, id, suggestion.StatusRejected, reviewerID, reason).Scan(&reviewedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	item.Status = suggestion.StatusRejected
	item.ReviewedBy = &reviewerID
	item.ReviewedAt = &reviewedAt
	item.RejectionReason = &reason

	return item, nil
}

func pendingForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*suggestion.Suggestion, error) {
	query := `select ` + suggestionColumns + ` ` + suggestionFrom + ` where s.id = $1 for update of s`

	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[suggestion.Suggestion])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, suggestions.ErrNotFound
		}
		return nil, err
	}

	if item.Status != suggestion.StatusPending {
		return nil, suggestions.ErrNotPending
	}

	return &item, nil
}
