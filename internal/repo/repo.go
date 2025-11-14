package repo

import (
	"context"

	e "app/internal/entity"
	"app/internal/repo/pgdb"
	rd "app/internal/repo/repodto"
	"app/pkg/postgres"
)

type Teams interface {
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
	CreateTeam(ctx context.Context, teamName string) (e.Team, error)
	DeleteUsersFromTeam(ctx context.Context, teamName string) error
}

type Users interface {
	UpsertBulk(ctx context.Context, user []e.User) error
	DeleteUsersByTeam(ctx context.Context, teamName string) error
	SetIsActive(ctx context.Context, userID string, isActive *bool) (e.User, error)
	GetUserByID(ctx context.Context, userID string) (e.User, error)
	GetActiveUsersTeam(ctx context.Context, teamName string, exIDs []string) ([]string, error)
}

type PullReq interface {
	GetPR(ctx context.Context, prID string) (e.PullRequest, error)
	CreatePR(ctx context.Context, in rd.CreatePRInput) (e.PullRequest, error)
	AssignReviewers(ctx context.Context, prID string, reviewers []string) ([]string, error)
	SetNeedMoreReviewrs(ctx context.Context, prID string, value bool) error
	ChangeReviewer(ctx context.Context, in rd.ChangeReviewerInput) error
	GetPRsByReviewer(ctx context.Context, uID string) ([]e.PullRequestShort, error)
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
