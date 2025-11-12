package pgdb

import (
	"context"

	e "app/internal/entity"
	"app/internal/repo/repoerrs"
	errutils "app/pkg/errors"
	"app/pkg/postgres"

	"github.com/jackc/pgx/v5"
)

type UsersRepo struct {
	*postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) Upsert(ctx context.Context, u e.User) error {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("user_id", "username", "team_name", "is_active").
		Values(u.UserID, u.Username, u.TeamName, u.IsActive).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	_, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return errutils.WrapPathErr(err)
	}

	return nil
}

func (r *UsersRepo) GetUsersByTeam(ctx context.Context, teamName string) ([]e.User, error) {
	sql, args, _ := r.Builder.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where("team_name = ?", teamName).
		OrderBy("user_id").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[e.User])

	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	if len(users) == 0 {
		return nil, repoerrs.ErrNotFound
	}

	return users, nil
}

func (r *UsersRepo) DeleteUsersByTeam(ctx context.Context, teamName string) error {
	sql, args, _ := r.Builder.
		Delete("users").
		Where("team_name = ?", teamName).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	cmdTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return errutils.WrapPathErr(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return repoerrs.ErrNoRowsDeleted
	}

	return nil
}
