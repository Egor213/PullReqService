package repo

import (
	"context"

	e "app/internal/entity"
	"app/internal/repo/pgdb"
	"app/pkg/postgres"
)

type Teams interface {
	CreateTeam(ctx context.Context, teamName string) (e.Team, error)
}

type Users interface {
	Upsert(ctx context.Context, user e.User) error
	GetUsersByTeam(ctx context.Context, teamName string) ([]e.User, error)
	DeleteUsersByTeam(ctx context.Context, teamName string) error
}

type PullReq interface {
	Temp(ctx context.Context) error
}

type Repositories struct {
	Users
	Teams
	PullReq
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Users:   pgdb.NewUsersRepo(pg),
		Teams:   pgdb.NewTeamsRepo(pg),
		PullReq: pgdb.NewPullReqRepo(pg),
	}
}
