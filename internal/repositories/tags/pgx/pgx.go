package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/repositories/tags"
	"github.com/Dmitro44/relper/internal/repositories/tags/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ tags.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func (p *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	query := `
select id, name, slug
from tags
where id = $1
`

	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	tag, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Tag])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &tag, nil
}

func (p *Pgx) Create(ctx context.Context, input tags.CreateInput) (*model.Tag, error) {
	query := `
insert into tags(name, slug)
values ($1, $2)
returning id, name, slug, created_at;
`

	rows, err := p.pool.Query(
		ctx, query,
		input.Name,
		input.Slug,
	)
	if err != nil {
		return nil, err
	}

	tag, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Tag])
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (p *Pgx) Update(ctx context.Context, id uuid.UUID, input tags.UpdateInput) (*model.Tag, error) {
	updateBuilder := sq.Update("tags").Where(squirrel.Eq{"id": id})

	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", input.Name)
	}
	if input.Slug != nil {
		updateBuilder = updateBuilder.Set("slug", input.Slug)
	}

	updateBuilder = updateBuilder.Suffix(`returning id, name, slug, created_at;`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	tag, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Tag])
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (p *Pgx) AddToPlace(ctx context.Context, placeID, tagID, userID uuid.UUID) error {
	query := `
insert into places_tags(place_id, tag_id, user_id)
values ($1, $2, $3)
`

	_, err := p.pool.Exec(ctx, query, placeID, tagID, userID)
	return err
}

func (p *Pgx) RemoveFromPlace(ctx context.Context, placeID, tagID uuid.UUID) error {
	query := `
delete from places_tags
where place_id = $1 and tag_id = $2
`

	_, err := p.pool.Exec(ctx, query, placeID, tagID)
	return err
}
