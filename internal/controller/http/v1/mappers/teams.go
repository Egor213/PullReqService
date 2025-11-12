package mappers

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/entity"
)

func ToEntityTeam(input httpdto.AddTeamInput) entity.Team {
	members := make([]entity.TeamMember, len(input.Members))
	for i, m := range input.Members {
		members[i] = entity.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return entity.Team{
		TeamName: input.TeamName,
		Members:  members,
	}
}

func ToAddTeamOutput(team entity.Team) httpdto.AddTeamOutput {
	members := make([]httpdto.TeamMemberInput, len(team.Members))
	for i, m := range team.Members {
		members[i] = httpdto.TeamMemberInput{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return httpdto.AddTeamOutput{
		Team: httpdto.AddTeamInput{
			TeamName: team.TeamName,
			Members:  members,
		},
	}
}
