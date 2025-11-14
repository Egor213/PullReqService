package pgdb

import (
	"app/pkg/postgres"
	"context"

	rd "app/internal/repo/dto"
	errutils "app/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type StatsRepo struct {
	*postgres.Postgres
}

func NewStatsRepo(pg *postgres.Postgres) *StatsRepo {
	return &StatsRepo{pg}
}

func (r *StatsRepo) GetReviewerStats(ctx context.Context) ([]rd.ReviewerStatsOutput, error) {
	sql, args, _ := r.Builder.
		Select("user_id", "COUNT(*) AS assignments").
		From("pr_reviewers").
		GroupBy("user_id").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}
	defer rows.Close()

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByName[rd.ReviewerStatsOutput])
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	return stats, nil
}

func (r *StatsRepo) GetPRStats(ctx context.Context) ([]rd.PRStatsOutput, error) {
	sql, args, _ := r.Builder.
		Select("pr_id", "COUNT(*) AS assignments").
		From("pr_reviewers").
		GroupBy("pr_id").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}
	defer rows.Close()

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByName[rd.PRStatsOutput])
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	return stats, nil
}
