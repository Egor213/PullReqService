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
