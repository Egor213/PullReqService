package service

import (
	"app/internal/repo"
)

type Services struct {
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{}
}
