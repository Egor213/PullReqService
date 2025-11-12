package repo

import (
	e "app/internal/entity"
	"app/pkg/postgres"
	"context"
)

type Teams interface {
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
	CreateTeam(ctx context.Context, teamName string) (e.Team, error)
}

type Users interface {
	Upsert(ctx context.Context, user e.User) error
	GetUsersByTeam(ctx context.Context, teamName string) (e.User, error)
	DeleteUsersByTeam(ctx context.Context, teamName string) error
}

type PullReq interface {
}

type Repositories struct {
	Users
	Teams
	PullReq
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{}
}
