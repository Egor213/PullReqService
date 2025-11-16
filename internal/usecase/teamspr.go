package usecase

import (
	"app/internal/service"
	serverrs "app/internal/service/errors"
	"context"
	"errors"
)

type TeamsPRUC struct {
	tService  service.Teams
	PRService service.PullReq
}

func NewTeamsPRUC(tService service.Teams, prService service.PullReq) *TeamsPRUC {
	return &TeamsPRUC{
		tService:  tService,
		PRService: prService,
	}
}

func (uc *TeamsPRUC) DeactivateTeamUsers(ctx context.Context, tName string) error {
	users, err := uc.tService.DeactivateTeamUsers(ctx, tName)
	if err != nil {
		return err
	}

	for _, u := range users {
		prs, err := uc.PRService.GetPRsByReviewer(ctx, u)
		if err != nil {
			return err
		}
		for _, pr := range prs {
			err := uc.PRService.DeleteReviewer(ctx, u, pr.PullReqID)
			if err != nil && !errors.Is(err, serverrs.ErrMergedPR) {
				return err
			}
		}
	}

	return nil
}
