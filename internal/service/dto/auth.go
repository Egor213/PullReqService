package servdto

import e "app/internal/entity"

type GenTokenInput struct {
	UserID string
	Role   e.Role
}
