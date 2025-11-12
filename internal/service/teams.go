package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/service/servdto"
	"context"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type TeamsService struct {
	teamsRepo repo.Teams
	usersRepo repo.Users
	trManager *manager.Manager
}

func NewTeamsService(tRepo repo.Teams, uRepo repo.Users, tr *manager.Manager) *TeamsService {
	return &TeamsService{
		usersRepo: uRepo,
		teamsRepo: tRepo,
		trManager: tr,
	}
}

func (s *TeamsService) CreateOrUpdateTeam(ctx context.Context, in e.Team) (e.Team, error) {
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		existingTeam, err := s.teamsRepo.GetTeam(ctx, in.TeamName)
		if err != nil && !errors.Is(err, repo.ErrTeamNotFound) {
			return err
		}

		if existingTeam != nil {
			if CompareMembers(existingTeam.Members, in.Members) {
				return ErrTeamWithUsersExists
			}

			if err := s.ReplaceTeamMembers(ctx, servdto.ReplaceMembersInput{
				TeamName: in.TeamName,
				Members:  in.Members,
			}); err != nil {
				return err
			}

			return nil
		}

		if err := s.teamsRepo.Create(ctx, in.TeamName); err != nil {
			return err
		}
		for _, m := range in.Members {
			if err := s.usersRepo.Upsert(ctx, m, in.TeamName); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return e.Team{}, err
	}

	return in, nil
}

func (s *TeamsService) ReplaceTeamMembers(ctx context.Context, in servdto.ReplaceMembersInput) error {
	return s.trManager.Do(ctx, func(ctx context.Context) error {
		currentMembers, err := s.usersRepo.GetUsersByTeam(ctx, in.TeamName)
		if err != nil {
			return err
		}

		if CompareMembers(currentMembers, in.Members) {
			return ErrTeamWithUsersExists
		}

		if err := s.usersRepo.DeleteUsersByTeam(ctx, in.TeamName); err != nil {
			return err
		}

		for _, m := range in.Members {
			if err := s.usersRepo.Upsert(ctx, m, in.TeamName); err != nil {
				return err
			}
		}

		return nil
	})
}
