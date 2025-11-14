package servdto

type ReplaceMembersInput struct {
	TeamName string
	Members  []TeamMemberDTO
}

type TeamMemberDTO struct {
	IsActive *bool
	UserID   string
	Username string
}

type CrOrUpTeamInput struct {
	TeamName string
	Members  []TeamMemberDTO
}
