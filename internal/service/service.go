package service

import (
	"context"

	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/service/servdto"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Teams interface {
	CreateOrUpdateTeam(ctx context.Context, in e.Team) (e.Team, error)
	ReplaceTeamMembers(ctx context.Context, in servdto.ReplaceMembersInput) error
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
}

type Users interface {
	SetIsActive(ctx context.Context, userID string, isActive *bool) (e.User, error)
}

type PullReq interface{}

type Services struct {
	Teams
	Users
	PullReq
}

type ServicesDependencies struct {
	Repos     *repo.Repositories
	TrManager *manager.Manager
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Teams: NewTeamsService(deps.Repos.Teams, deps.Repos.Users, deps.TrManager),
		Users: NewUsersService(deps.Repos.Users),
	}
}
