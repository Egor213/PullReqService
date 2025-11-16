package httpdto

type TeamMemberInput struct {
	UserID   string `json:"user_id" validate:"required,max=100"`
	Username string `json:"username" validate:"required,max=100"`
	IsActive *bool  `json:"is_active" validate:"required"`
}

type AddTeamInput struct {
	TeamName string            `json:"team_name" validate:"required,max=100"`
	Members  []TeamMemberInput `json:"members" validate:"required,min=1,dive"`
}

type AddTeamOutput struct {
	Team AddTeamInput `json:"team"`
}

type GetTeamInput struct {
	TeamName string `query:"team_name" validate:"required,max=100"`
}

type GetTeamOutput struct {
	TeamName string            `json:"team_name"`
	Members  []TeamMemberInput `json:"members"`
}

type DeactivateTeamInput struct {
	TeamName string `json:"team_name" validate:"required,max=100"`
}

type DeactivateTeamOutput struct {
	TeamName string `json:"team_name"`
}
