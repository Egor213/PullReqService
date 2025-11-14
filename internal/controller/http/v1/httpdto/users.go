package httpdto

import e "app/internal/entity"

type SetIsActiveInput struct {
	UserID   string `json:"user_id" validate:"required,max=100"`
	IsActive *bool  `json:"is_active" validate:"required"`
}

type SetIsActiveOutput struct {
	User e.User `json:"user"`
}
