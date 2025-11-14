package service

import (
	"context"
	"errors"

	e "app/internal/entity"
	"app/internal/repo"
	re "app/internal/repo/errors"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
	smap "app/internal/service/mappers"

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

func (s *TeamsService) CreateOrUpdateTeam(ctx context.Context, in sd.CrOrUpTeamInput) (e.Team, error) {
	members := smap.TeamMemberDTOToMember(in.Members)
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		team, err := s.teamsRepo.GetTeam(ctx, in.TeamName)

		if err != nil {
			if !errors.Is(err, re.ErrNotFound) {
				return err
			}

			_, err = s.teamsRepo.CreateTeam(ctx, in.TeamName)
			if err != nil {
				return err
			}

			err := s.usersRepo.UpsertBulk(ctx, smap.TeamMemberDTOToUser(in.Members, in.TeamName))
			if err != nil {
				return err
			}

			return nil
		}

		if CompareMembers(team.Members, members) {
			return se.ErrTeamWithUsersExists
		}

		err = s.ReplaceTeamMembers(ctx, sd.ReplaceMembersInput(in))
		if err != nil && !errors.Is(err, re.ErrNoRowsDeleted) {
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

func (s *TeamsService) ReplaceTeamMembers(ctx context.Context, in sd.ReplaceMembersInput) error {
	return s.trManager.Do(ctx, func(ctx context.Context) error {
		err := s.teamsRepo.DeleteUsersFromTeam(ctx, in.TeamName)

		if err != nil && !errors.Is(err, re.ErrNoRowsDeleted) {
			return err
		}

		err = s.usersRepo.UpsertBulk(ctx, smap.TeamMemberDTOToUser(in.Members, in.TeamName))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *TeamsService) GetTeam(ctx context.Context, teamName string) (e.Team, error) {
	team, err := s.teamsRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) {
			return e.Team{}, se.ErrNotFoundTeam
		}
		return e.Team{}, err
	}

	return team, nil
}
