package usecase

import (
	"app/internal/entity"
	"app/internal/service"
	servdto "app/internal/service/dto"
	"app/internal/usecase/dto"
	"context"
)

type UsersPRUC struct {
	uService  service.Users
	PRService service.PullReq
}

func NewUsersPRUC(uService service.Users, prService service.PullReq) *UsersPRUC {
	return &UsersPRUC{
		uService:  uService,
		PRService: prService,
	}
}

func (uc *UsersPRUC) SetIsActiveAndReassignPRs(ctx context.Context, in dto.ActiveAndReassugnInput) (entity.User, error) {
	user, err := uc.uService.SetIsActive(ctx, servdto.SetIsActiveInput{
		UserID:   in.UserID,
		IsActive: in.IsActive,
	})
	if err != nil {
		return entity.User{}, err
	}

	prs, err := uc.PRService.GetPRsByReviewer(ctx, user.UserID)
	if err != nil {
		return entity.User{}, err
	}

	isForce := true
	for _, pr := range prs {
		_, err = uc.PRService.ReassignReviewer(ctx, servdto.ReassignReviewerInput{
			PullReqID: pr.PullReqID,
			RevID:     user.UserID,
			Force:     &isForce,
		})
		if err != nil {
			return entity.User{}, err
		}
	}

	return user, err
}
