package service

import (
	"app/internal/repo"
	"context"
	"errors"

	e "app/internal/entity"

	re "app/internal/repo/errors"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
	smap "app/internal/service/mappers"
	errutils "app/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type TeamsService struct {
	teamsRepo repo.Teams
	usersRepo repo.Users
	trManager TRManager
}

func NewTeamsService(tRepo repo.Teams, uRepo repo.Users, tr TRManager) *TeamsService {
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
			log.Error(errutils.WrapPathErr(err))
			if !errors.Is(err, re.ErrNotFound) {
				return se.ErrCannotGetTeam
			}

			_, err = s.teamsRepo.CreateTeam(ctx, in.TeamName)
			if err != nil {
				log.Error(errutils.WrapPathErr(err))
				if errors.Is(err, re.ErrAlreadyExists) {
					return se.ErrTeamWExists
				}
				return se.ErrCannotCreateTeam
			}

			err := s.usersRepo.UpsetBulk(ctx, smap.TeamMemberDTOToUser(in.Members, in.TeamName))
			if err != nil {
				log.Error(errutils.WrapPathErr(err))
				return se.ErrCannotUpsetUsers
			}

			return nil
		}

		if CompareMembers(team.Members, members) {
			return se.ErrTeamWithUsersExists
		}

		err = s.ReplaceTeamMembers(ctx, sd.ReplaceMembersInput(in))
		if err != nil {
			log.Error(errutils.WrapPathErr(err))
			return err
		}

		return nil
	})
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
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

		if err != nil && !errors.Is(err, re.ErrNotFound) {
			log.Error(errutils.WrapPathErr(err))
			return se.ErrCannotDelUsersFromTeam
		}

		err = s.usersRepo.UpsetBulk(ctx, smap.TeamMemberDTOToUser(in.Members, in.TeamName))
		if err != nil {
			log.Error(errutils.WrapPathErr(err))
			return se.ErrCannotUpsetUsers
		}

		return nil
	})
}

func (s *TeamsService) GetTeam(ctx context.Context, teamName string) (e.Team, error) {
	team, err := s.teamsRepo.GetTeam(ctx, teamName)
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		return e.Team{}, se.HandleRepoNotFound(err, se.ErrNotFoundTeam, se.ErrCannotGetTeam)
	}

	return team, nil
}

func (s *TeamsService) DeactivateTeamUsers(ctx context.Context, teamName string) ([]string, error) {
	users, err := s.teamsRepo.DeactivateTeamUsers(ctx, teamName)
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		return nil, se.ErrCannotDeactivateTeam
	}
	return users, nil
}
