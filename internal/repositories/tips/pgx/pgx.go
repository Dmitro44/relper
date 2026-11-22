package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/entities/tip"
	"github.com/Dmitro44/relper/internal/repositories/tips"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ tips.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewTipsRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	query := `
select id, location_id, place_id, user_id, category, title, content, created_at
from tips
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tip.Tip])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tips.ErrNotFound
		}
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) ListByPlace(ctx context.Context, placeID uuid.UUID) ([]tip.Tip, error) {
	query := `
select id, location_id, place_id, user_id, category, title, content, created_at
from tips
where place_id = $1
order by created_at desc
`

	rows, err := p.pool.Query(ctx, query, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[tip.Tip])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]tip.Tip, 0)
	}

	return items, nil
}

func (p *Pgx) Create(ctx context.Context, input tips.CreateInput) (*tip.Tip, error) {
	query := `
insert into tips(location_id, place_id, user_id, category, title, content)
values ($1, $2, $3, $4, $5, $6)
returning id, location_id, place_id, user_id, category, title, content, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.LocationID,
		input.PlaceID,
		input.UserID,
		input.Category,
		input.Title,
		input.Content,
	)
	if err != nil {
		return nil, err
	}

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tip.Tip])
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input tips.UpdateInput) (*tip.Tip, error) {
	updateBuilder := sq.Update("tips").Where(squirrel.Eq{"id": id})

	if !admin {
		updateBuilder = updateBuilder.Where(squirrel.Eq{"user_id": userID})
	}

	if input.LocationID != nil {
		updateBuilder = updateBuilder.Set("location_id", input.LocationID).Set("place_id", nil)
	}
	if input.PlaceID != nil {
		updateBuilder = updateBuilder.Set("place_id", input.PlaceID).Set("location_id", nil)
	}
	if input.Category != nil {
		updateBuilder = updateBuilder.Set("category", input.Category)
	}
	if input.Title != nil {
		updateBuilder = updateBuilder.Set("title", input.Title)
	}
	if input.Content != nil {
		updateBuilder = updateBuilder.Set("content", input.Content)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, location_id, place_id, user_id, category, title, content, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tip.Tip])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tips.ErrNotFound
		}
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error {
	builder := sq.Delete("tips").Where(squirrel.Eq{"id": id})

	if !admin {
		builder = builder.Where(squirrel.Eq{"user_id": userID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	tag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return tips.ErrNotFound
	}

	return nil
}
