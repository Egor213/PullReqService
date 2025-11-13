package pgdb

import (
	"context"
	"errors"

	e "app/internal/entity"
	"app/internal/repo/repoerrs"
	errutils "app/pkg/errors"
	"app/pkg/postgres"

	sq "github.com/Masterminds/squirrel"

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

func (r *UsersRepo) SetIsActive(ctx context.Context, userID string, isActive *bool) (e.User, error) {
	sql, args, _ := r.Builder.
		Update("users").
		Set("is_active", isActive).
		Where("user_id = ?", userID).
		Suffix("RETURNING user_id, username, team_name, is_active").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	row := conn.QueryRow(ctx, sql, args...)

	var user e.User
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.User{}, repoerrs.ErrNotFound
		}
		return e.User{}, errutils.WrapPathErr(err)
	}

	return user, nil
}

func (r *UsersRepo) GetUserByID(ctx context.Context, userID string) (e.User, error) {
	sql, args, _ := r.Builder.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where("user_id = ?", userID).
		Limit(1).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	row := conn.QueryRow(ctx, sql, args...)

	var user e.User
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.User{}, repoerrs.ErrNotFound
		}
		return e.User{}, errutils.WrapPathErr(err)
	}

	return user, nil
}

func (r *UsersRepo) GetActiveUsersTeam(ctx context.Context, teamName string, exIDs []string) ([]string, error) {
	builder := r.Builder.
		Select("user_id").
		From("users").
		Where("team_name = ?", teamName).
		Where("is_active = true")

	if len(exIDs) > 0 {
		builder = builder.Where(sq.NotEq{"user_id": exIDs})
	}

	sql, args, _ := builder.ToSql()
	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, errutils.WrapPathErr(err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errutils.WrapPathErr(err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}
