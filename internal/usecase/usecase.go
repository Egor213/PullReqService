package usecase

import (
	"app/internal/entity"
	"app/internal/service"
	"app/internal/usecase/dto"
	"context"
)

type UsersPRUseCase interface {
	SetIsActiveAndReassignPRs(ctx context.Context, in dto.ActiveAndReassugnInput) (entity.User, error)
}

type TeamsPRUseCase interface {
	DeactivateTeamUsers(ctx context.Context, tName string) error
}

type UseCases struct {
	UsersPRUseCase
	TeamsPRUseCase
}

type UseCasesDependencies struct {
	Servs *service.Services
}

func NewUseCases(dep UseCasesDependencies) *UseCases {
	return &UseCases{
		UsersPRUseCase: NewUsersPRUC(dep.Servs.Users, dep.Servs.PullReq),
		TeamsPRUseCase: NewTeamsPRUC(dep.Servs.Teams, dep.Servs.PullReq),
	}
}
