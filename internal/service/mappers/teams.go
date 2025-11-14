package servmappers

import (
	e "app/internal/entity"
	servdto "app/internal/service/dto"
)

func TeamMemberDTOToUser(t []servdto.TeamMemberDTO, team string) []e.User {
	u := make([]e.User, 0, len(t))
	for _, m := range t {
		u = append(u, e.User{
			IsActive: m.IsActive,
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: team,
		})
	}
	return u
}

func TeamMemberDTOToMember(t []servdto.TeamMemberDTO) []e.TeamMember {
	tm := make([]e.TeamMember, 0, len(t))
	for _, m := range t {
		tm = append(tm, e.TeamMember{
			IsActive: m.IsActive,
			UserID:   m.UserID,
			Username: m.Username,
		})
	}
	return tm
}
