package httpmappers

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/entity"
	"app/internal/service/servdto"
)

func ToCrOrUpTeamInput(input httpdto.AddTeamInput) servdto.CrOrUpTeamInput {
	members := make([]servdto.TeamMemberDTO, len(input.Members))
	for i, m := range input.Members {
		members[i] = servdto.TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return servdto.CrOrUpTeamInput{
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
