package service

import (
	"context"
	"errors"

	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repoerrs"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"
)

type UsersService struct {
	usersRepo repo.Users
}

func NewUsersService(uRepo repo.Users) *UsersService {
	return &UsersService{
		usersRepo: uRepo,
	}
}

func (s *UsersService) SetIsActive(ctx context.Context, in servdto.SetIsActiveInput) (e.User, error) {
	user, err := s.usersRepo.SetIsActive(ctx, in.UserID, in.IsActive)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return e.User{}, serverrs.ErrUserNotFound
		}
		return e.User{}, serverrs.ErrCannotSetParam
	}

	return user, nil
}
