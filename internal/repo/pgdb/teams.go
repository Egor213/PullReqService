package pgdb

import (
	"context"
	"errors"

	errutils "app/pkg/errors"
	"app/pkg/postgres"

	e "app/internal/entity"
	entitymappers "app/internal/entity/mappers"
	repoerrs "app/internal/repo/repoerrs"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TeamsRepo struct {
	*postgres.Postgres
}

func NewTeamsRepo(pg *postgres.Postgres) *TeamsRepo {
	return &TeamsRepo{pg}
}

func (r *TeamsRepo) CreateTeam(ctx context.Context, teamName string) (e.Team, error) {
	sql, args, _ := r.Builder.
		Insert("teams").
		Columns("team_name").
		Values(teamName).
		Suffix("RETURNING team_name").
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	var team e.Team
	if err := conn.QueryRow(ctx, sql, args...).Scan(&team.TeamName); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Team{}, repoerrs.ErrNotFound
		} else if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return e.Team{}, repoerrs.ErrAlreadyExists
		}
		return e.Team{}, errutils.WrapPathErr(err)
	}

	return team, nil
}

func (r *TeamsRepo) GetTeam(ctx context.Context, teamName string) (e.Team, error) {
	sql, args, _ := r.Builder.
		Select("team_name").
		From("teams").
		Where("team_name = ?", teamName).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	row := conn.QueryRow(ctx, sql, args...)

	var team e.Team
	if err := row.Scan(&team.TeamName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Team{}, repoerrs.ErrNotFound
		}
		return e.Team{}, errutils.WrapPathErr(err)
	}

	sql, args, _ = r.Builder.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where("team_name = ?", teamName).
		OrderBy("user_id").
		ToSql()

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return e.Team{}, errutils.WrapPathErr(err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[e.User])
	if err != nil {
		return e.Team{}, errutils.WrapPathErr(err)
	}

	team.Members = entitymappers.UsersToTeamMembers(users)

	return team, nil
}

func (r *TeamsRepo) DeleteUsersFromTeam(ctx context.Context, teamName string) error {
	sql, args, _ := r.Builder.
		Update("users").
		Set("team_name", nil).
		Where("team_name = ?", teamName).
		ToSql()

	conn := r.CtxGetter.DefaultTrOrDB(ctx, r.Pool)
	_, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return errutils.WrapPathErr(err)
	}

	return nil
}
