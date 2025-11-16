package httpdto

type SetIsActiveInput struct {
	UserID   string `json:"user_id" validate:"required,max=100"`
	IsActive *bool  `json:"is_active" validate:"required"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive *bool  `json:"is_active"`
}

type SetIsActiveOutput struct {
	User UserDTO `json:"user"`
}
