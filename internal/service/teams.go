package service

import (
	"context"
	"errors"

	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repoerrs"
	servmappers "app/internal/service/mappers"
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

func (s *TeamsService) CreateOrUpdateTeam(ctx context.Context, in servdto.CrOrUpTeamInput) (e.Team, error) {
	members := servmappers.TeamMemberDTOToMember(in.Members)
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		team, err := s.teamsRepo.GetTeam(ctx, in.TeamName)

		if err != nil {
			if !errors.Is(err, repoerrs.ErrNotFound) {
				return err
			}

			_, err = s.teamsRepo.CreateTeam(ctx, in.TeamName)
			if err != nil {
				return err
			}

			err := s.usersRepo.UpsertBulk(ctx, servmappers.TeamMemberDTOToUser(in.Members, in.TeamName))
			if err != nil {
				return err
			}

			return nil
		}

		if CompareMembers(team.Members, members) {
			return serverrs.ErrTeamWithUsersExists
		}

		err = s.ReplaceTeamMembers(ctx, servdto.ReplaceMembersInput(in))
		if err != nil && !errors.Is(err, repoerrs.ErrNoRowsDeleted) {
			return err
		}

		return nil
	})

	if err != nil {
		return e.Team{}, err
	}

	return e.Team{
		TeamName: in.TeamName,
		Members:  members,
	}, nil
}

func (s *TeamsService) ReplaceTeamMembers(ctx context.Context, in servdto.ReplaceMembersInput) error {
	return s.trManager.Do(ctx, func(ctx context.Context) error {
		err := s.teamsRepo.DeleteUsersFromTeam(ctx, in.TeamName)

		if err != nil && !errors.Is(err, repoerrs.ErrNoRowsDeleted) {
			return err
		}

		err = s.usersRepo.UpsertBulk(ctx, servmappers.TeamMemberDTOToUser(in.Members, in.TeamName))
		if err != nil {
			return err
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
