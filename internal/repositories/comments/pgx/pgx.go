package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/entities/comment"
	"github.com/Dmitro44/relper/internal/repositories/comments"
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

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	query := `
select id, place_id, tip_id, user_id, content, created_at
from comments
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	comment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, comments.ErrNotFound
		}
		return nil, err
	}

	return &comment, nil
}

func (p *Pgx) ListByPlace(ctx context.Context, placeID uuid.UUID) ([]comment.WithAuthor, error) {
	query := `
select c.id, c.user_id, u.nickname, c.content, c.created_at
from comments c
join users u on u.id = c.user_id
where c.place_id = $1
order by c.created_at asc
`

	rows, err := p.pool.Query(ctx, query, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[comment.WithAuthor])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]comment.WithAuthor, 0)
	}

	return items, nil
}

func (p *Pgx) ListByTip(ctx context.Context, tipID uuid.UUID) ([]comment.WithAuthor, error) {
	query := `
select c.id, c.user_id, u.nickname, c.content, c.created_at
from comments c
join users u on u.id = c.user_id
where c.tip_id = $1
order by c.created_at asc
`

	rows, err := p.pool.Query(ctx, query, tipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[comment.WithAuthor])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]comment.WithAuthor, 0)
	}

	return items, nil
}

func (p *Pgx) Create(ctx context.Context, input comments.CreateInput) (*comment.Comment, error) {
	query := `
insert into comments(place_id, tip_id, user_id, content)
values ($1, $2, $3, $4)
returning id, place_id, tip_id, user_id, content, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.PlaceID,
		input.TipID,
		input.UserID,
		input.Content,
	)
	if err != nil {
		return nil, err
	}

	comment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (p *Pgx) Update(ctx context.Context, id, userID uuid.UUID, admin bool, input comments.UpdateInput) (*comment.Comment, error) {
	updateBuilder := sq.Update("comments").Where(squirrel.Eq{"id": id})

	if !admin {
		updateBuilder = updateBuilder.Where(squirrel.Eq{"user_id": userID})
	}

	if input.Content != nil {
		updateBuilder = updateBuilder.Set("content", input.Content)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, place_id, tip_id, user_id, content, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	comment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, comments.ErrNotFound
		}
		return nil, err
	}

	return &comment, nil
}

func (p *Pgx) Delete(ctx context.Context, id, userID uuid.UUID, admin bool) error {
	builder := sq.Delete("comments").Where(squirrel.Eq{"id": id})

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
		return comments.ErrNotFound
	}

	return nil
}
