package pgx

import (
	"context"

	auditentity "github.com/Dmitro44/relper/internal/entities/audit"
	auditrepo "github.com/Dmitro44/relper/internal/repositories/audit"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ auditrepo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewAuditRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{pool: pool}
}

func (p *Pgx) Insert(ctx context.Context, input auditrepo.CreateInput) error {
	details := input.Details
	if len(details) == 0 {
		details = []byte("{}")
	}

	_, err := p.pool.Exec(ctx, `
insert into audit_log (source, action, entity_type, entity_id, actor_id, details)
values ($1, $2, $3, $4, $5, $6);
`, input.Source, input.Action, input.EntityType, input.EntityID, input.ActorID, details)

	return err
}

func (p *Pgx) List(ctx context.Context, filter auditentity.Filter) ([]auditentity.Event, error) {
	query := sq.Select(
		"a.id",
		"a.source",
		"a.action",
		"a.entity_type",
		"a.entity_id",
		"a.actor_id",
		"coalesce(u.nickname, '') as actor_nickname",
		"a.details",
		"a.created_at",
	).
		From("audit_log a").
		LeftJoin("users u ON u.id = a.actor_id").
		OrderBy("a.created_at DESC").
		Limit(uint64(filter.Limit))

	if filter.Source != "" {
		query = query.Where(squirrel.Eq{"a.source": filter.Source})
	}
	if filter.Action != "" {
		query = query.Where(squirrel.Eq{"a.action": filter.Action})
	}
	if filter.EntityType != "" {
		query = query.Where(squirrel.Eq{"a.entity_type": filter.EntityType})
	}
	if filter.ActorID != nil {
		query = query.Where(squirrel.Eq{"a.actor_id": *filter.ActorID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[auditentity.Event])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]auditentity.Event, 0)
	}

	return items, nil
}
