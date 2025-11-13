package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repodto"
	"app/internal/repo/repoerrs"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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
			return serverrs.ErrCreatePR
		}
		return nil
	})
	if err != nil {
		return e.PullRequest{}, err
	}
	return pr, nil
}
