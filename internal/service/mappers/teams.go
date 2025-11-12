package servmappers

import (
	e "app/internal/entity"
)

func UsersToTeamMembers(users []e.User) []e.TeamMember {
	members := make([]e.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, e.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}
	return members
}

func TeamMemberToUser(member e.TeamMember, teamName string) e.User {
	return e.User{
		UserID:   member.UserID,
		Username: member.Username,
		TeamName: teamName,
		IsActive: member.IsActive,
	}
}
