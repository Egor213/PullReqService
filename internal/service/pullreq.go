package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	rd "app/internal/repo/repodto"
	re "app/internal/repo/repoerrs"
	sd "app/internal/service/servdto"
	se "app/internal/service/serverrs"
	errutils "app/pkg/errors"
	"context"
	"errors"
	"math/rand/v2"
	"slices"

	"github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type PullReqService struct {
	prRepo    repo.PullReq
	usersRepo repo.Users
	trManager *manager.Manager
}

func NewPullReqService(prRepo repo.PullReq, uRepo repo.Users, tr *manager.Manager) *PullReqService {
	return &PullReqService{
		prRepo:    prRepo,
		usersRepo: uRepo,
		trManager: tr,
	}
}

func (s *PullReqService) CreatePR(ctx context.Context, in sd.CreatePRInput) (e.PullRequest, error) {
	var pr e.PullRequest
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		user, err := s.usersRepo.GetUserByID(ctx, in.AuthorID)
		if err != nil {
			if errors.Is(err, re.ErrNotFound) {
				return se.ErrUserNotFound
			}
			return se.ErrCannotGetUser
		}

		if user.IsActive == nil || !*user.IsActive {
			return se.ErrInactiveCreator
		}

		pr, err = s.prRepo.CreatePR(ctx, rd.CreatePRInput{
			PullReqID: in.PullReqID,
			NamePR:    in.NamePR,
			AuthorID:  in.AuthorID,
			Status:    e.StatusOpen,
		})
		if err != nil {
			if errors.Is(err, re.ErrAlreadyExists) {
				return se.ErrPRExists
			}
			log.Error(errutils.WrapPathErr(err))
			return se.ErrCreatePR
		}

		pr.Reviewers, err = s.AssignReviewers(ctx, sd.AssignReviewersInput{
			PullReqID:    pr.PullReqID,
			AuthorTeam:   user.TeamName,
			ExcludeUsers: []string{in.AuthorID},
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return e.PullRequest{}, err
	}
	return pr, nil
}

func (s *PullReqService) AssignReviewers(ctx context.Context, in sd.AssignReviewersInput) ([]string, error) {
	const reviewersCount = 2
	var out []string
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		users, err := s.usersRepo.GetActiveUsersTeam(ctx, in.AuthorTeam, in.ExcludeUsers)
		if err != nil {
			return se.ErrCannotGetUser
		}

		if len(users) == 0 {
			if err := s.prRepo.SetNeedMoreReviewrs(ctx, in.PullReqID, true); err != nil {
				return err
			}
			return nil
		}

		rand.Shuffle(len(users), func(i, j int) {
			users[i], users[j] = users[j], users[i]
		})

		count := min(reviewersCount, len(users))
		selectedReviewers := users[:count]
		out, err = s.prRepo.AssignReviewers(ctx, in.PullReqID, selectedReviewers)
		if err != nil {
			return err
		}

		needMore := count < reviewersCount
		if err := s.prRepo.SetNeedMoreReviewrs(ctx, in.PullReqID, needMore); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PullReqService) GetPR(ctx context.Context, prID string) (e.PullRequest, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) {
			return e.PullRequest{}, se.ErrNotFoundPR
		}
		// TODO: везде такое поставить где неизвестная ошибка
		log.Error(errutils.WrapPathErr(err))
		return e.PullRequest{}, se.ErrCannotGetPR
	}
	return pr, nil
}

func (s *PullReqService) ReassignReviewer(ctx context.Context, in sd.ReassignReviewerInput) (sd.ReassignReviewerOutput, error) {
	trOp := pgxv5.MustSettings(
		settings.Must(),
		pgxv5.WithTxOptions(pgx.TxOptions{IsoLevel: pgx.RepeatableRead}),
	)
	var out sd.ReassignReviewerOutput
	err := s.trManager.DoWithSettings(ctx, trOp, func(ctx context.Context) error {
		pr, err := s.prRepo.GetPR(ctx, in.PullReqID)

		if err != nil {
			if errors.Is(err, re.ErrNotFound) {
				return se.ErrNotFoundPR
			}
			log.Error(errutils.WrapPathErr(err))
			return se.ErrCannotGetPR
		}

		if len(pr.Reviewers) == 0 {
			return se.ErrNotFoundReviewers
		}
		log.Info(pr.Reviewers)
		if !slices.Contains(pr.Reviewers, in.RevID) {
			return se.ErrReviewerNotAssigned
		}

		if pr.Status == e.StatusMerged {
			return se.ErrMergedPR
		}

		exIDs := []string{pr.AuthorID}
		exIDs = append(exIDs, pr.Reviewers...)

		user, err := s.usersRepo.GetUserByID(ctx, pr.AuthorID)
		if err != nil {
			if errors.Is(err, re.ErrNotFound) {
				return se.ErrUserNotFound
			}
			return se.ErrCannotGetUser
		}

		users, err := s.usersRepo.GetActiveUsersTeam(ctx, user.TeamName, exIDs)
		if err != nil {
			return se.ErrCannotGetUser
		}

		if len(users) == 0 {
			return se.ErrNoAvailableReviewers
		}

		NewRevID := users[rand.IntN(len(users))]
		err = s.prRepo.ChangeReviewer(ctx, rd.ChangeReviewerInput{
			PullReqID:   in.PullReqID,
			NewReviewer: NewRevID,
			OldReviewer: in.RevID,
		})

		if err != nil {
			if errors.Is(err, re.ErrNotFound) {
				return se.ErrNotFoundReviewers
			}
			log.Error(errutils.WrapPathErr(err))
			return err
		}

		out.NewRevID = NewRevID
		id := slices.Index(pr.Reviewers, in.RevID)
		pr.Reviewers[id] = NewRevID
		out.PullReq = pr

		return nil
	})
	if err != nil {
		return sd.ReassignReviewerOutput{}, err
	}
	return out, nil
}

func (s *PullReqService) GetPRsByReviewer(ctx context.Context, uID string) ([]e.PullRequestShort, error) {
	_, err := s.usersRepo.GetUserByID(ctx, uID)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) {
			return nil, se.ErrUserNotFound
		}
		return nil, se.ErrCannotGetUser
	}

	prs, err := s.prRepo.GetPRsByReviewer(ctx, uID)
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		return nil, se.ErrCannotGetPR
	}
	return prs, nil
}
