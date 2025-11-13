package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repodto"
	"app/internal/repo/repoerrs"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"
	errutils "app/pkg/errors"
	"context"
	"errors"
	"math/rand/v2"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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

func (s *PullReqService) CreatePR(ctx context.Context, in servdto.CreatePRInput) (e.PullRequest, error) {
	var pr e.PullRequest
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		user, err := s.usersRepo.GetUserByID(ctx, in.AuthorID)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return serverrs.ErrUserNotFound
			}
			return serverrs.ErrCannotGetUser
		}

		if user.IsActive == nil || !*user.IsActive {
			return serverrs.ErrInactiveCreator
		}

		pr, err = s.prRepo.CreatePR(ctx, repodto.CreatePRInput{
			PullReqID: in.PullReqID,
			NamePR:    in.NamePR,
			AuthorID:  in.AuthorID,
			Status:    e.StatusOpen,
		})
		if err != nil {
			if errors.Is(err, repoerrs.ErrAlreadyExists) {
				return serverrs.ErrPRExists
			}
			log.Error(errutils.WrapPathErr(err))
			return serverrs.ErrCreatePR
		}

		err = s.AssignReviewers(ctx, servdto.AssignReviewersInput{
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

func (s *PullReqService) AssignReviewers(ctx context.Context, in servdto.AssignReviewersInput) error {
	const reviewersCount = 2

	return s.trManager.Do(ctx, func(ctx context.Context) error {
		users, err := s.usersRepo.GetActiveUsersTeam(ctx, in.AuthorTeam, in.ExcludeUsers)
		if err != nil {
			return serverrs.ErrCannotGetUser
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

		if _, err := s.prRepo.AssignReviewers(ctx, in.PullReqID, selectedReviewers); err != nil {
			return err
		}

		needMore := count < reviewersCount
		if err := s.prRepo.SetNeedMoreReviewrs(ctx, in.PullReqID, needMore); err != nil {
			return err
		}

		return nil
	})
}

func (s *PullReqService) GetPR(ctx context.Context, prID string) (e.PullRequest, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return e.PullRequest{}, serverrs.ErrNotFoundPR
		}
		// TODO: везде такое поставить где неизвестная ошибка
		log.Error(errutils.WrapPathErr(err))
		return e.PullRequest{}, serverrs.ErrCannotGetPR
	}
	return pr, nil
}
