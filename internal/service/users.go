package service

import (
	"context"
	"errors"

	e "app/internal/entity"
	"app/internal/repo"
	re "app/internal/repo/errors"
	sd "app/internal/service/servdto"
	se "app/internal/service/serverrs"
)

type UsersService struct {
	usersRepo repo.Users
}

func NewUsersService(uRepo repo.Users) *UsersService {
	return &UsersService{
		usersRepo: uRepo,
	}
}

func (s *UsersService) SetIsActive(ctx context.Context, in sd.SetIsActiveInput) (e.User, error) {
	user, err := s.usersRepo.SetIsActive(ctx, in.UserID, in.IsActive)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) {
			return e.User{}, se.ErrUserNotFound
		}
		return e.User{}, se.ErrCannotSetParam
	}

	return user, nil
}
