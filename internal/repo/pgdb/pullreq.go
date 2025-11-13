package pgdb

import (
	e "app/internal/entity"
	"app/internal/repo/repodto"
	repoerrs "app/internal/repo/repoerrs"
	"app/pkg/postgres"
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PullReqRepo struct {
	*postgres.Postgres
}

func NewPullReqRepo(pg *postgres.Postgres) *PullReqRepo {
	return &PullReqRepo{pg}
}

func (r *PullReqRepo) CreatePR(ctx context.Context, in repodto.CreatePRInput) (e.PullRequest, error) {
	sql, args, _ := r.Builder.
		Insert("prs").
		Columns("pr_id", "title", "author_id", "status").
		Values(in.PullReqID, in.NamePR, in.AuthorID, in.Status).
		Suffix("RETURNING pr_id, title, author_id, status, created_at, merged_at").
		ToSql()
	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)

	var pr e.PullRequest
	err := conn.QueryRow(ctx, sql, args...).Scan(
		&pr.PullReqID,
		&pr.NamePR,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return e.PullRequest{}, repoerrs.ErrAlreadyExists
		case pgerrcode.ForeignKeyViolation:
			return e.PullRequest{}, repoerrs.ErrNotFound

		}
	}

	return pr, nil
}
