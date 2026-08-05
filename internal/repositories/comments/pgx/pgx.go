package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/repositories/comments"
	"github.com/Dmitro44/relper/internal/repositories/comments/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ comments.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewCommentsRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	query := `
select id, target_type, target_id, user_id, content, created_at
from comments
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	comment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Comment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &comment, nil
}

func (p *Pgx) Create(ctx context.Context, input comments.CreateInput) (*model.Comment, error) {
	query := `
insert into comments(target_type, target_id, user_id, content)
values ($1, $2, $3, $4)
returning id, target_type, target_id, user_id, content, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.TargetType,
		input.TargetID,
		input.UserID,
		input.Content,
	)
	if err != nil {
		return nil, err
	}

	comment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Comment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &comment, nil
}

func (p *Pgx) Update(ctx context.Context, id uuid.UUID, input comments.UpdateInput) (*model.Comment, error) {
	updateBuilder := sq.Update("comments").Where(squirrel.Eq{"id": id})

	if input.TargetType != nil {
		updateBuilder = updateBuilder.Set("target_type", input.TargetType)
	}
	if input.TargetID != nil {
		updateBuilder = updateBuilder.Set("target_id", input.TargetID)
	}
	if input.UserID != nil {
		updateBuilder = updateBuilder.Set("user_id", input.UserID)
	}
	if input.Content != nil {
		updateBuilder = updateBuilder.Set("content", input.Content)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, target_type, target_id, user_id, content, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	comment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Comment])
	if err != nil {
		return nil, err
	}

	return &comment, nil
}
