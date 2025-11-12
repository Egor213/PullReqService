package pgdb

import (
	errutils "app/pkg/errors"
	"app/pkg/postgres"
	"context"
	"errors"

	e "app/internal/entity"
	repoerrs "app/internal/repo/repoerrs"

	"github.com/jackc/pgerrcode"
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
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return e.Team{}, repoerrs.ErrAlreadyExists
		}
		return e.Team{}, errutils.WrapPathErr(err)
	}

	return team, nil
}
