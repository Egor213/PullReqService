package servdto

import e "app/internal/entity"

type ReplaceMembersInput struct {
	TeamName string
	Members  []e.TeamMember
}
