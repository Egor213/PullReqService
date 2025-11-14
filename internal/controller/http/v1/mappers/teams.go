package mapper

import (
	hd "app/internal/controller/http/v1/httpdto"
	e "app/internal/entity"
	sd "app/internal/service/servdto"
)

func ToTeamMemberDTO(m e.TeamMember) sd.TeamMemberDTO {
	return sd.TeamMemberDTO{
		UserID:   m.UserID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

func ToTeamMemberDTOs(members []e.TeamMember) []sd.TeamMemberDTO {
	dtos := make([]sd.TeamMemberDTO, len(members))
	for i, m := range members {
		dtos[i] = ToTeamMemberDTO(m)
	}
	return dtos
}

func ToTeamMemberInput(m e.TeamMember) hd.TeamMemberInput {
	return hd.TeamMemberInput{
		UserID:   m.UserID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

func ToTeamMemberInputs(members []e.TeamMember) []hd.TeamMemberInput {
	inp := make([]hd.TeamMemberInput, len(members))
	for i, m := range members {
		inp[i] = ToTeamMemberInput(m)
	}
	return inp
}

func ToTeamMemberDTOFromInput(m hd.TeamMemberInput) sd.TeamMemberDTO {
	return sd.TeamMemberDTO{
		UserID:   m.UserID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

func ToTeamMemberDTOsFromInput(members []hd.TeamMemberInput) []sd.TeamMemberDTO {
	dtos := make([]sd.TeamMemberDTO, len(members))
	for i, m := range members {
		dtos[i] = ToTeamMemberDTOFromInput(m)
	}
	return dtos
}

func ToCrOrUpTeamInput(input hd.AddTeamInput) sd.CrOrUpTeamInput {
	return sd.CrOrUpTeamInput{
		TeamName: input.TeamName,
		Members:  ToTeamMemberDTOsFromInput(input.Members),
	}
}

func ToAddTeamOutput(team e.Team) hd.AddTeamOutput {
	return hd.AddTeamOutput{
		Team: hd.AddTeamInput{
			TeamName: team.TeamName,
			Members:  ToTeamMemberInputs(team.Members),
		},
	}
}

func ToPullRequestShortDTO(pr e.PullRequestShort) hd.PullRequestShortDTO {
	return hd.PullRequestShortDTO{
		PullReqID: pr.PullReqID,
		NamePR:    pr.NamePR,
		AuthorID:  pr.AuthorID,
		Status:    string(pr.Status),
	}
}

func ToPullRequestShortDTOs(prs []e.PullRequestShort) []hd.PullRequestShortDTO {
	dtos := make([]hd.PullRequestShortDTO, 0, len(prs))
	for _, pr := range prs {
		dtos = append(dtos, ToPullRequestShortDTO(pr))
	}
	return dtos
}

func ToGetReviewOutput(userID string, prs []e.PullRequestShort) hd.GetReviewOutput {
	return hd.GetReviewOutput{
		UserID:  userID,
		PullReq: ToPullRequestShortDTOs(prs),
	}
}
