package pgx

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/rating"
	"github.com/Dmitro44/relper/internal/repositories/ratings"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ ratings.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewRatingsRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) Upsert(ctx context.Context, userID, placeID uuid.UUID, stars int16) (*rating.PlaceRating, error) {
	query := `
insert into place_ratings (user_id, place_id, stars)
values ($1, $2, $3)
on conflict (user_id, place_id)
do update set stars = excluded.stars, updated_at = now()
returning id, user_id, place_id, stars, created_at, updated_at;
`

	rows, err := p.pool.Query(ctx, query, userID, placeID, stars)
	if err != nil {
		return nil, err
	}

	r, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rating.PlaceRating])
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (p *Pgx) Delete(ctx context.Context, userID, placeID uuid.UUID) error {
	query := `
delete from place_ratings
where user_id = $1 and place_id = $2;
`

	_, err := p.pool.Exec(ctx, query, userID, placeID)

	return err
}
