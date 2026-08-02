package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/repositories/places"
	"github.com/Dmitro44/relper/internal/repositories/places/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ places.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewPlacesRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) Create(ctx context.Context, input places.CreateInput) (*model.Place, error) {
	query := `
insert into places(location_id, name, category, address, phone, website, coordinates, description, rating, upvotes)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning id, location_id, name, category, address, phone, website, coordinates, description, rating, upvotes, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.LocationID,
		input.Name,
		input.Category,
		input.Address,
		input.Phone,
		input.WebsiteURL,
		input.Coordinates,
		input.Description,
		input.Rating,
		input.Upvotes,
	)
	if err != nil {
		return nil, err
	}

	place, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Place])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &place, nil
}

func (p *Pgx) Update(ctx context.Context, id uuid.UUID, input places.UpdateInput) (*model.Place, error) {
	updateBuilder := sq.Update("places").Where(squirrel.Eq{"id": id})

	if input.LocationID != nil {
		updateBuilder = updateBuilder.Set("location_id", input.LocationID)
	}
	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", input.Name)
	}
	if input.Category != nil {
		updateBuilder = updateBuilder.Set("category", input.Category)
	}
	if input.Address != nil {
		updateBuilder = updateBuilder.Set("address", input.Address)
	}
	if input.Phone != nil {
		updateBuilder = updateBuilder.Set("phone", input.Phone)
	}
	if input.WebsiteURL != nil {
		updateBuilder = updateBuilder.Set("website", input.WebsiteURL)
	}
	if input.Coordinates != nil {
		updateBuilder = updateBuilder.Set("coordinates", input.Coordinates)
	}
	if input.Description != nil {
		updateBuilder = updateBuilder.Set("description", input.Description)
	}
	if input.Rating != nil {
		updateBuilder = updateBuilder.Set("rating", input.Rating)
	}
	if input.Upvotes != nil {
		updateBuilder = updateBuilder.Set("upvotes", input.Upvotes)
	}
	if input.Upvotes != nil {
		updateBuilder = updateBuilder.Set("upvotes", input.Upvotes)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, location_id, name, category, address, phone, website, coordinates, description, rating, upvotes, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	place, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Place])
	if err != nil {
		return nil, err
	}

	return &place, nil
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*model.Place, error) {
	query := `
select id, location_id, name, category, address, phone, website, coordinates, description, rating, upvotes, created_at 
from places
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	place, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Place])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &place, nil
}
