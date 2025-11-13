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
	trManager *manager.Manager
}

func NewPullReqService(prRepo repo.PullReq, tr *manager.Manager) *PullReqService {
	return &PullReqService{
		prRepo:    prRepo,
		trManager: tr,
	}
}

func (s *PullReqService) CreatePR(ctx context.Context, in servdto.CreatePRInput) (e.PullRequest, error) {
	pr, err := s.prRepo.CreatePR(ctx, repodto.CreatePRInput{
		PullReqID: in.PullReqID,
		NamePR:    in.NamePR,
		AuthorID:  in.AuthorID,
		Status:    e.StatusOpen,
	})
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return e.PullRequest{}, serverrs.ErrPRExists
		} else if errors.Is(err, repoerrs.ErrNotFound) {
			return e.PullRequest{}, serverrs.ErrCannotGetUser
		}
		return e.PullRequest{}, serverrs.ErrCreatePR
	}
	return pr, nil
}
