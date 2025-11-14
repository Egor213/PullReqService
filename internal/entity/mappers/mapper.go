package entitymappers

import e "app/internal/entity"

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

func TeamMembersToUsers(members []e.TeamMember, teamName string) []e.User {
	users := make([]e.User, 0, len(members))
	for _, m := range members {
		users = append(users, e.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: teamName,
			IsActive: m.IsActive,
		})
	}
	return users
}
