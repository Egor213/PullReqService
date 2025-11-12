package repo

import (
	"app/pkg/postgres"
)

type Repositories struct {
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{}
}
