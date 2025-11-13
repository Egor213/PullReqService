package service

import (
	"context"
	"errors"

	e "app/internal/entity"
	entitymappers "app/internal/entity/mappers"
	"app/internal/repo"
	"app/internal/repo/repoerrs"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"

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
		team, err := s.teamsRepo.GetTeam(ctx, in.TeamName)

		if err != nil {
			if !errors.Is(err, repoerrs.ErrNotFound) {
				return err
			}

			if _, err := s.teamsRepo.CreateTeam(ctx, in.TeamName); err != nil {
				return err
			}

			for _, m := range in.Members {
				if err := s.usersRepo.Upsert(ctx, entitymappers.TeamMemberToUser(m, in.TeamName)); err != nil {
					return err
				}
			}

			return nil
		}

		if CompareMembers(team.Members, in.Members) {
			return serverrs.ErrTeamWithUsersExists
		}

		if err := s.ReplaceTeamMembers(ctx, servdto.ReplaceMembersInput{
			TeamName: in.TeamName,
			Members:  in.Members,
		}); err != nil && !errors.Is(err, repoerrs.ErrNoRowsDeleted) {
			return err
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
		err := s.teamsRepo.DeleteUsersFromTeam(ctx, in.TeamName)

		if err != nil && !errors.Is(err, repoerrs.ErrNoRowsDeleted) {
			return err
		}

		for _, m := range in.Members {
			err := s.usersRepo.Upsert(ctx, entitymappers.TeamMemberToUser(m, in.TeamName))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *TeamsService) GetTeam(ctx context.Context, teamName string) (e.Team, error) {
	team, err := s.teamsRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return e.Team{}, serverrs.ErrNotFoundTeam
		}
		return e.Team{}, err
	}

	return team, nil
}
