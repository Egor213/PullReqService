package pgdb

import (
	"app/pkg/postgres"
	"context"
)

type PullReqRepo struct {
	*postgres.Postgres
}

func NewPullReqRepo(pg *postgres.Postgres) *PullReqRepo {
	return &PullReqRepo{pg}
}

func (r *PullReqRepo) Temp(ctx context.Context) error {
	return nil
}
