package repo

import (
	"context"
	"time"

	e "app/internal/entity"
	rd "app/internal/repo/dto"
	"app/internal/repo/pgdb"
	"app/pkg/postgres"
)

type Teams interface {
	GetTeam(ctx context.Context, teamName string) (e.Team, error)
	CreateTeam(ctx context.Context, teamName string) (e.Team, error)
	DeleteUsersFromTeam(ctx context.Context, teamName string) error
}

type Users interface {
	UpsetBulk(ctx context.Context, user []e.User) error
	DeleteUsersByTeam(ctx context.Context, teamName string) error
	SetIsActive(ctx context.Context, userID string, isActive *bool) (e.User, error)
	GetUserByID(ctx context.Context, userID string) (e.User, error)
	GetActiveUsersTeam(ctx context.Context, teamName string, exIDs []string) ([]string, error)
}

type PullReq interface {
	GetPR(ctx context.Context, prID string) (e.PullRequest, error)
	CreatePR(ctx context.Context, in rd.CreatePRInput) (e.PullRequest, error)
	AssignReviewers(ctx context.Context, prID string, reviewers []string) ([]string, error)
	SetNeedMoreReviewers(ctx context.Context, prID string, value bool) error
	ChangeReviewer(ctx context.Context, in rd.ChangeReviewerInput) error
	GetPRsByReviewer(ctx context.Context, uID string) ([]e.PullRequestShort, error)
	MergePR(ctx context.Context, prID string) (*time.Time, error)
}

type Stats interface {
	GetReviewerStats(ctx context.Context) ([]rd.ReviewerStatsOutput, error)
	GetPRStats(ctx context.Context) ([]rd.PRStatsOutput, error)
}

type Repositories struct {
	Users
	Teams
	PullReq
	Stats
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Users:   pgdb.NewUsersRepo(pg),
		Teams:   pgdb.NewTeamsRepo(pg),
		PullReq: pgdb.NewPullReqRepo(pg),
		Stats:   pgdb.NewStatsRepo(pg),
	}
}
