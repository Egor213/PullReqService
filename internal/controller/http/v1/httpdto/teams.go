package httpdto

import "app/internal/entity"

type AddTeamInput struct {
	TeamName string              `json:"team_name" validate:"required"`
	Members  []entity.TeamMember `json:"members" validate:"required,dive"`
}

type AddTeamOutput struct {
	Team AddTeamInput `json:"team"`
}
