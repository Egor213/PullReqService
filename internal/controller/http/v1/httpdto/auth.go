package httpdto

import e "app/internal/entity"

// TODO: ПЕРЕДАТЬ ВСЕ ПОД  INPUT - OUTPUT DTO даже если уже готовые сущности

type LoginInput struct {
	UserID string `json:"user_id" validate:"required"`
	Role   e.Role `json:"role" validate:"required,oneof=admin user"`
}

type LoginOutput struct {
	AccessToken string `json:"access_token"`
}
