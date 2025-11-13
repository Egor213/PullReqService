package service

import (
	"context"
	"errors"

	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repoerrs"
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

func (s *UsersService) SetIsActive(ctx context.Context, userID string, isActive *bool) (e.User, error) {
	user, err := s.usersRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return e.User{}, serverrs.ErrUserNotFound
		}
		return e.User{}, err
	}

	return user, nil
}
