package pgx

import (
	"context"

	"github.com/google/uuid"

	rolesrepo "github.com/Dmitro44/relper/internal/repositories/roles"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ rolesrepo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewRolesRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) Add(ctx context.Context, userID uuid.UUID, codes []string) error {
	query := `
insert into user_roles (user_id, role_id)
select $1, id from roles where code = any($2::text[])
on conflict do nothing
`

	_, err := p.pool.Exec(ctx, query, userID, codes)

	return err
}

func (p *Pgx) GetByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
select r.code
from user_roles ur
join roles r on r.id = ur.role_id
where ur.user_id = $1
`

	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (p *Pgx) ReplaceForUser(ctx context.Context, userID uuid.UUID, codes []string) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	roleIDs, err := p.roleIDsByCode(ctx, tx, codes)
	if err != nil {
		return err
	}
	query := `
delete from user_roles
where user_id = $1
`

	if _, err = tx.Exec(ctx, query, userID); err != nil {
		return err
	}

	query = `
insert into user_roles (user_id, role_id)
values ($1, $2)
on conflict do nothing
`
	for _, roleID := range roleIDs {
		if _, err = tx.Exec(ctx, query, userID, roleID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *Pgx) roleIDsByCode(ctx context.Context, tx pgx.Tx, codes []string) ([]uuid.UUID, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	query := `
select id from roles 
where code = any($1::text[])
`

	rows, err := tx.Query(ctx, query, codes)
	if err != nil {
		return nil, err
	}

	roleIDs, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, err
	}

	if len(roleIDs) != len(codes) {
		return nil, rolesrepo.ErrUnknownRole
	}

	return roleIDs, nil
}
