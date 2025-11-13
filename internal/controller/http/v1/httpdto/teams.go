package httpdto

type TeamMemberInput struct {
	IsActive *bool  `json:"is_active" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type AddTeamInput struct {
	TeamName string            `json:"team_name" validate:"required"`
	Members  []TeamMemberInput `json:"members" validate:"required,min=1,dive"`
}

type AddTeamOutput struct {
	Team AddTeamInput `json:"team"`
}

type GetTeamInput struct {
	TeamName string `query:"team_name" validate:"required"`
}
