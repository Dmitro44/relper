package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/repositories/tips"
	"github.com/Dmitro44/relper/internal/repositories/tips/model"
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

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*model.Tip, error) {
	query := `
select id, location_id, user_id, category, title, content, upvotes, created_at
from tips
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Tip])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) Create(ctx context.Context, input tips.CreateInput) (*model.Tip, error) {
	query := `
insert into tips(location_id, user_id, category, title, content, upvotes)
values ($1, $2, $3, $4, $5, $6)
returning id, location_id, user_id, category, title, content, upvotes, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.LocationID,
		input.UserID,
		input.Category,
		input.Title,
		input.Content,
		input.Upvotes,
	)
	if err != nil {
		return nil, err
	}

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Tip])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) Update(ctx context.Context, id uuid.UUID, input tips.UpdateInput) (*model.Tip, error) {
	updateBuilder := sq.Update("tips").Where(squirrel.Eq{"id": id})

	if input.LocationID != nil {
		updateBuilder = updateBuilder.Set("location_id", input.LocationID)
	}
	if input.UserID != nil {
		updateBuilder = updateBuilder.Set("user_id", input.UserID)
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
	if input.Upvotes != nil {
		updateBuilder = updateBuilder.Set("upvotes", input.Upvotes)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, location_id, user_id, category, title, content, upvotes, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Tip])
	if err != nil {
		return nil, err
	}

	return &tip, nil
}
