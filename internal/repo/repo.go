package repo

import (
	"context"

	e "app/internal/entity"
	"app/internal/repo/pgdb"
	"app/internal/repo/repodto"
	"app/pkg/postgres"
)

type Teams interface {
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
	CreateTeam(ctx context.Context, teamName string) (e.Team, error)
	DeleteUsersFromTeam(ctx context.Context, teamName string) error
}

type Users interface {
	Upsert(ctx context.Context, user e.User) error
	DeleteUsersByTeam(ctx context.Context, teamName string) error
	SetIsActive(ctx context.Context, userID string, isActive *bool) (e.User, error)
	GetUserByID(ctx context.Context, userID string) (e.User, error)
}

type PullReq interface {
	CreatePR(ctx context.Context, in repodto.CreatePRInput) (e.PullRequest, error)
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
