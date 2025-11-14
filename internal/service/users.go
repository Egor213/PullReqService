package service

import (
	"context"

	e "app/internal/entity"
	"app/internal/repo"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
	errutils "app/pkg/errors"

	log "github.com/sirupsen/logrus"
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
		log.Error(errutils.WrapPathErr(err))
		return e.User{}, se.HandleRepoNotFound(err, se.ErrNotFoundUser, se.ErrCannotSetParam)
	}

	return user, nil
}
