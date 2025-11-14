package service

import (
	"context"
	"time"

	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/service/servdto"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Teams interface {
	CreateOrUpdateTeam(ctx context.Context, in servdto.CrOrUpTeamInput) (e.Team, error)
	ReplaceTeamMembers(ctx context.Context, in servdto.ReplaceMembersInput) error
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
}

type Users interface {
	SetIsActive(ctx context.Context, in servdto.SetIsActiveInput) (e.User, error)
}

type PullReq interface {
	CreatePR(ctx context.Context, in servdto.CreatePRInput) (e.PullRequest, error)
	AssignReviewers(ctx context.Context, in servdto.AssignReviewersInput) error
	GetPR(ctx context.Context, prID string) (e.PullRequest, error)
}

type Auth interface {
	GenerateToken(ctx context.Context, in servdto.GenTokenInput) (string, error)
	ParseToken(accessToken string) (e.ParsedToken, error)
}

type Services struct {
	Teams
	Users
	PullReq
	Auth
}

type ServicesDependencies struct {
	Repos     *repo.Repositories
	TrManager *manager.Manager

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Teams:   NewTeamsService(deps.Repos.Teams, deps.Repos.Users, deps.TrManager),
		Auth:    NewAuthService(deps.Repos.Users, deps.SignKey, deps.TokenTTL),
		Users:   NewUsersService(deps.Repos.Users),
		PullReq: NewPullReqService(deps.Repos.PullReq, deps.Repos.Users, deps.TrManager),
	}
}
