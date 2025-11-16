package service

import (
	"app/internal/repo"
	"context"
	"time"

	e "app/internal/entity"

	sd "app/internal/service/dto"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type TRManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
	DoWithSettings(ctx context.Context, s trm.Settings, fn func(ctx context.Context) error) (err error)
}

type Teams interface {
	CreateOrUpdateTeam(ctx context.Context, in sd.CrOrUpTeamInput) (e.Team, error)
	ReplaceTeamMembers(ctx context.Context, in sd.ReplaceMembersInput) error
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
	DeactivateTeamUsers(ctx context.Context, teamName string) ([]string, error)
}

type Users interface {
	SetIsActive(ctx context.Context, in sd.SetIsActiveInput) (e.User, error)
}

type PullReq interface {
	CreatePR(ctx context.Context, in sd.CreatePRInput) (e.PullRequest, error)
	AssignReviewers(ctx context.Context, in sd.AssignReviewersInput) ([]string, error)
	GetPR(ctx context.Context, prID string) (e.PullRequest, error)
	ReassignReviewer(ctx context.Context, in sd.ReassignReviewerInput) (sd.ReassignReviewerOutput, error)
	GetPRsByReviewer(ctx context.Context, uID string) ([]e.PullRequestShort, error)
	MergePR(ctx context.Context, prID string) (e.PullRequest, error)
	DeleteReviewer(ctx context.Context, uID string, prID string) error
}

type Auth interface {
	GenerateToken(ctx context.Context, in sd.GenTokenInput) (string, error)
	ParseToken(accessToken string) (e.ParsedToken, error)
}

type Stats interface {
	GetStats(ctx context.Context) (sd.GetStatsOutput, error)
}

type Services struct {
	Teams
	Users
	PullReq
	Auth
	Stats
}

type ServicesDependencies struct {
	Repos     *repo.Repositories
	TrManager TRManager

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Teams:   NewTeamsService(deps.Repos.Teams, deps.Repos.Users, deps.TrManager),
		Auth:    NewAuthService(deps.Repos.Users, deps.SignKey, deps.TokenTTL),
		Users:   NewUsersService(deps.Repos.Users),
		PullReq: NewPullReqService(deps.Repos.PullReq, deps.Repos.Users, deps.TrManager),
		Stats:   NewStatsService(deps.Repos.Stats),
	}
}
