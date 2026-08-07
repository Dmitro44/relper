package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/entities/location"
	"github.com/Dmitro44/relper/internal/repositories/locations"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ locations.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewLocationsRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) Create(ctx context.Context, input locations.CreateInput) (*location.Location, error) {
	query := `
insert into locations(country, city, coordinates)
values ($1, $2, $3)
returning id, country, city, coordinates, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.Country,
		input.City,
		input.Coordinates,
	)
	if err != nil {
		return nil, err
	}

	location, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[location.Location])
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (p *Pgx) Update(ctx context.Context, id uuid.UUID, input locations.UpdateInput) (*location.Location, error) {
	updateBuilder := sq.Update("locations").Where(squirrel.Eq{"id": id})

	if input.Country != nil {
		updateBuilder = updateBuilder.Set("country", input.Country)
	}
	if input.City != nil {
		updateBuilder = updateBuilder.Set("city", input.City)
	}
	if input.Coordinates != nil {
		updateBuilder = updateBuilder.Set("coordinates", input.Coordinates)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, country, city, coordinates, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	location, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[location.Location])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, locations.ErrNotFound
		}
		return nil, err
	}

	return &location, nil
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error) {
	query := `
select id, country, city, coordinates
from locations
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	location, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[location.Location])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, locations.ErrNotFound
		}
		return nil, err
	}

	return &location, nil
}

func (p *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
delete from locations
where id = $1
`

	tag, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return locations.ErrNotFound
	}

	return nil
}
