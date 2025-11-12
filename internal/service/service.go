package service

import (
	"app/internal/repo"
	"app/internal/service/servdto"
	"context"
)

type Teams interface {
	AddTeam(ctx context.Context, in servdto.AddTeamInput) error
}

type PullReq interface {
}

type Services struct {
	Teams
	PullReq
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{}
}
