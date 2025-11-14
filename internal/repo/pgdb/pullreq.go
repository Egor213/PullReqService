package pgdb

import (
	e "app/internal/entity"
	rd "app/internal/repo/repodto"
	repoerrs "app/internal/repo/repoerrs"
	errutils "app/pkg/errors"
	"app/pkg/postgres"
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PullReqRepo struct {
	*postgres.Postgres
}

func NewPullReqRepo(pg *postgres.Postgres) *PullReqRepo {
	return &PullReqRepo{pg}
}

func (r *PullReqRepo) CreatePR(ctx context.Context, in rd.CreatePRInput) (e.PullRequest, error) {
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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return e.PullRequest{}, repoerrs.ErrAlreadyExists
			case pgerrcode.ForeignKeyViolation:
				return e.PullRequest{}, repoerrs.ErrNotFound
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return e.PullRequest{}, repoerrs.ErrNotFound
		}

		return e.PullRequest{}, errutils.WrapPathErr(err)
	}

	return pr, nil
}

func (r *PullReqRepo) GetPR(ctx context.Context, prID string) (e.PullRequest, error) {
	sql, args, _ := r.Builder.
		Select(
			"prs.pr_id", "prs.title", "prs.author_id", "prs.status", "prs.need_more_reviewers", "prs.created_at", "prs.merged_at",
			"array_agg(pr_reviewers.user_id) AS assigned_reviewers",
		).
		From("prs").
		LeftJoin("pr_reviewers ON pr_reviewers.pr_id = prs.pr_id").
		Where("prs.pr_id = ?", prID).
		GroupBy("prs.pr_id, prs.title, prs.author_id, prs.status, prs.created_at, prs.merged_at").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)

	var pr e.PullRequest
	var assignedReviewers []*string
	err := conn.QueryRow(ctx, sql, args...).Scan(
		&pr.PullReqID,
		&pr.NamePR,
		&pr.AuthorID,
		&pr.Status,
		&pr.NeedMoreReviewers,
		&pr.CreatedAt,
		&pr.MergedAt,
		&assignedReviewers,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.PullRequest{}, repoerrs.ErrNotFound
		}
		return e.PullRequest{}, errutils.WrapPathErr(err)
	}

	pr.Reviewers = make([]string, 0, len(assignedReviewers))
	for _, s := range assignedReviewers {
		if s != nil {
			pr.Reviewers = append(pr.Reviewers, *s)
		}
	}
	return pr, nil
}

func (r *PullReqRepo) AssignReviewers(ctx context.Context, prID string, reviewers []string) ([]string, error) {
	if len(reviewers) == 0 {
		return nil, nil
	}

	builder := r.Builder.
		Insert("pr_reviewers").
		Columns("pr_id", "user_id")

	for _, userID := range reviewers {
		builder = builder.Values(prID, userID)
	}

	builder = builder.Suffix("ON CONFLICT (pr_id, user_id) DO NOTHING")

	sql, args, _ := builder.ToSql()
	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)

	_, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	return reviewers, nil
}

func (r *PullReqRepo) SetNeedMoreReviewrs(ctx context.Context, prID string, value bool) error {
	sql, args, _ := r.Builder.
		Update("prs").
		Set("need_more_reviewers", value).
		Where("pr_id = ?", prID).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	cmdTag, err := conn.Exec(ctx, sql, args...)

	if err != nil {
		return errutils.WrapPathErr(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return repoerrs.ErrNotFound
	}

	return nil
}

func (r *PullReqRepo) ChangeReviewer(ctx context.Context, in rd.ChangeReviewerInput) error {
	sql, args, _ := r.Builder.
		Update("pr_reviewers").
		Set("user_id", in.NewReviewer).
		Where("pr_id = ? AND user_id = ?", in.PullReqID, in.OldReviewer).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	cmdTag, err := conn.Exec(ctx, sql, args...)

	if err != nil {
		return errutils.WrapPathErr(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return repoerrs.ErrNotFound
	}
	return nil
}

func (r *PullReqRepo) GetPRsByReviewer(ctx context.Context, uID string) ([]e.PullRequestShort, error) {
	sql, args, _ := r.Builder.
		Select("prs.pr_id", "prs.title", "prs.author_id", "prs.status").
		From("prs").
		Join("pr_reviewers ON pr_reviewers.pr_id = prs.pr_id").
		Where("pr_reviewers.user_id = ?", uID).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}
	defer rows.Close()

	prs, err := pgx.CollectRows(rows, pgx.RowToStructByName[e.PullRequestShort])
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	return prs, nil
}
