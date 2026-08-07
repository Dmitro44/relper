package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/Dmitro44/relper/internal/repositories/users"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ users.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewUsersRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) Create(ctx context.Context, input users.CreateInput) (*user.User, error) {
	createQuery := `
insert into users (email, nickname, habits, contacts)
values ($1, $2, $3, $4)
returning id, email, nickname, habits, contacts, created_at;
`
	rows, err := p.pool.Query(
		ctx, createQuery,
		input.Email,
		input.Nickname,
		input.Habits,
		input.Contacts,
	)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Pgx) Update(ctx context.Context, id uuid.UUID, input users.UpdateInput) (*user.User, error) {
	updateBuilder := sq.Update("users").Where(squirrel.Eq{"id": id})

	if input.Email != nil {
		updateBuilder = updateBuilder.Set("email", input.Email)
	}
	if input.Nickname != nil {
		updateBuilder = updateBuilder.Set("nickname", input.Nickname)
	}
	if input.Habits != nil {
		updateBuilder = updateBuilder.Set("habits", input.Habits)
	}
	if input.Contacts != nil {
		updateBuilder = updateBuilder.Set("contacts", input.Contacts)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, email, nickname, habits, contacts, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
select id, email, nickname, habits, contacts, created_at
from users
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (p *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
delete from users
where id = $1
`

	tag, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return users.ErrNotFound
	}

	return nil
}
