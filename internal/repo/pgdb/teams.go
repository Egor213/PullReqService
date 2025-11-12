package pgdb

import (
	e "app/internal/entity"
	"app/pkg/postgres"
	"context"
)

type TeamsRepo struct {
	*postgres.Postgres
}

func NewTeamsRepo(pg *postgres.Postgres) *TeamsRepo {
	return &TeamsRepo{pg}
}

func (r *TeamsRepo) GetTeam(ctx context.Context, teamName string) (e.Team, error) {
	return e.Team{}, nil
}

func (r *TeamsRepo) CreateTeam(ctx context.Context, teamName string) (e.Team, error) {
	return e.Team{}, nil
}
