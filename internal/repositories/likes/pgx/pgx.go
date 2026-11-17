package pgx

import (
	"context"

	"github.com/Dmitro44/relper/internal/repositories/likes"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ likes.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func NewLikesRepo(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool: pool,
	}
}

func (p *Pgx) LikeTip(ctx context.Context, userID, tipID uuid.UUID) error {
	query := `
insert into likes (user_id, tip_id)
values ($1, $2)
on conflict (user_id, tip_id) where tip_id is not null do nothing;
`

	_, err := p.pool.Exec(ctx, query, userID, tipID)

	return err
}

func (p *Pgx) LikeComment(ctx context.Context, userID, commentID uuid.UUID) error {
	query := `
insert into likes (user_id, comment_id)
values ($1, $2)
on conflict (user_id, comment_id) where comment_id is not null do nothing;
`

	_, err := p.pool.Exec(ctx, query, userID, commentID)

	return err
}

func (p *Pgx) UnlikeTip(ctx context.Context, userID, tipID uuid.UUID) error {
	return p.unlike(ctx, userID, tipID, "tip_id")
}

func (p *Pgx) UnlikeComment(ctx context.Context, userID, commentID uuid.UUID) error {
	return p.unlike(ctx, userID, commentID, "comment_id")
}

func (p *Pgx) unlike(ctx context.Context, userID, targetID uuid.UUID, column string) error {
	query := `delete from likes where user_id = $1 and ` + column + ` = $2;`

	_, err := p.pool.Exec(ctx, query, userID, targetID)

	return err
}

func (p *Pgx) LikedTip(ctx context.Context, userID, tipID uuid.UUID) (bool, error) {
	return p.liked(ctx, userID, tipID, "tip_id")
}

func (p *Pgx) LikedComment(ctx context.Context, userID, commentID uuid.UUID) (bool, error) {
	return p.liked(ctx, userID, commentID, "comment_id")
}

func (p *Pgx) liked(ctx context.Context, userID, targetID uuid.UUID, column string) (bool, error) {
	query := `select exists(select 1 from likes where user_id = $1 and ` + column + ` = $2);`

	var liked bool
	if err := p.pool.QueryRow(ctx, query, userID, targetID).Scan(&liked); err != nil {
		return false, err
	}

	return liked, nil
}

func (p *Pgx) CountTip(ctx context.Context, tipID uuid.UUID) (int64, error) {
	return p.count(ctx, tipID, "tip_id")
}

func (p *Pgx) CountComment(ctx context.Context, commentID uuid.UUID) (int64, error) {
	return p.count(ctx, commentID, "comment_id")
}

func (p *Pgx) count(ctx context.Context, targetID uuid.UUID, column string) (int64, error) {
	query := `select count(*)::bigint from likes where ` + column + ` = $1;`

	var count int64
	if err := p.pool.QueryRow(ctx, query, targetID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
