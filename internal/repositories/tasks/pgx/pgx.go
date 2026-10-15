package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/entities/task"
	"github.com/Dmitro44/relper/internal/repositories/tasks"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ tasks.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewTasksRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	query := `
select id, location_id, user_id, title, description, place_id, done, created_at
from tasks
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	task, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[task.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tasks.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (p *Pgx) Create(ctx context.Context, input tasks.CreateInput) (*task.Task, error) {
	query := `
insert into tasks(location_id, user_id, title, description, place_id, done)
values ($1, $2, $3, $4, $5, $6)
returning id, location_id, user_id, title, description, place_id, done, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.LocationID,
		input.UserID,
		input.Title,
		input.Description,
		input.PlaceID,
		input.Done,
	)
	if err != nil {
		return nil, err
	}

	task, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[task.Task])
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (p *Pgx) Update(ctx context.Context, id, userID uuid.UUID, input tasks.UpdateInput) (*task.Task, error) {
	updateBuilder := sq.Update("tasks").Where(squirrel.Eq{"id": id, "user_id": userID})

	if input.LocationID != nil {
		updateBuilder = updateBuilder.Set("location_id", input.LocationID)
	}
	if input.Title != nil {
		updateBuilder = updateBuilder.Set("title", input.Title)
	}
	if input.Description != nil {
		updateBuilder = updateBuilder.Set("description", input.Description)
	}
	if input.PlaceID != nil {
		updateBuilder = updateBuilder.Set("place_id", input.PlaceID)
	}
	if input.Done != nil {
		updateBuilder = updateBuilder.Set("done", input.Done)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, location_id, user_id, title, description, place_id, done, created_at`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	task, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[task.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tasks.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (p *Pgx) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `
delete from tasks
where id = $1 and user_id = $2
`

	tag, err := p.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return tasks.ErrNotFound
	}

	return nil
}
