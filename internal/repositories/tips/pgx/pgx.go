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
select id, location_id, user_id, category, title, content, upvotes, created_at
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

func (p *Pgx) Create(ctx context.Context, input tips.CreateInput) (*tip.Tip, error) {
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

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tip.Tip])
	if err != nil {
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) Update(ctx context.Context, id, userID uuid.UUID, input tips.UpdateInput) (*tip.Tip, error) {
	updateBuilder := sq.Update("tips").Where(squirrel.Eq{"id": id, "user_id": userID})

	if input.LocationID != nil {
		updateBuilder = updateBuilder.Set("location_id", input.LocationID)
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

	tip, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[tip.Tip])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tips.ErrNotFound
		}
		return nil, err
	}

	return &tip, nil
}

func (p *Pgx) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `
delete from tips
where id = $1 and user_id = $2
`

	tag, err := p.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return tips.ErrNotFound
	}

	return nil
}

func (p *Pgx) Upvote(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	query := `
update tips
set upvotes = upvotes + 1
where id = $1
returning id, location_id, user_id, category, title, content, upvotes, created_at;
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

func (p *Pgx) RemoveUpvote(ctx context.Context, id uuid.UUID) (*tip.Tip, error) {
	query := `
update tips
set upvotes = GREATEST(upvotes - 1, 0)
where id = $1
returning id, location_id, user_id, category, title, content, upvotes, created_at;
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
