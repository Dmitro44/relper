package pgx

import (
	"context"
	"errors"

	"github.com/Dmitro44/relper/internal/entities/refreshtoken"
	"github.com/Dmitro44/relper/internal/repositories/refreshtokens"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ refreshtokens.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) Create(ctx context.Context, input refreshtokens.CreateInput) error {
	query := `
insert into refresh_tokens(user_id, token_hash, expires_at)
values ($1, $2, $3)
returning id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by
`

	rows, err := p.pool.Query(
		ctx, query,
		input.UserID,
		input.TokenHash,
		input.ExpiresAt,
	)
	if err != nil {
		return err
	}

	_, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[refreshtoken.RefreshToken])
	if err != nil {
		return err
	}

	return nil
}

func (p *Pgx) GetByHash(ctx context.Context, tokenHash string) (*refreshtoken.RefreshToken, error) {
	query := `
select id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by
from refresh_tokens
where token_hash = $1
`

	rows, err := p.pool.Query(ctx, query, tokenHash)
	if err != nil {
		return nil, err
	}

	token, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[refreshtoken.RefreshToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, refreshtokens.ErrNotFound
		}
		return nil, err
	}

	return &token, nil
}

func (p *Pgx) RevokeByHash(ctx context.Context, tokenHash string) error {
	query := `
update refresh_tokens set revoked_at = now()
where token_hash = $1
`

	_, err := p.pool.Exec(ctx, query, tokenHash)

	return err
}

func (p *Pgx) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	query := `
update refresh_tokens set revoked_at = now()
where user_id = $1
	and revoked_at is null
	and expires_at > now()
`

	_, err := p.pool.Exec(ctx, query, userID)

	return err
}
