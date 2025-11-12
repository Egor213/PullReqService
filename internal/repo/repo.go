package repo

import (
	"app/pkg/postgres"
)

type Teams interface {
}

type Users interface {
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
