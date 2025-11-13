package serverrs

import (
	"errors"
)

// TODO: добавь нормальные ошибки в методах сервисов, а то че там просто err надо через эти ошибки
var (
	ErrTeamWithUsersExists = errors.New("team with such users already exists")
	ErrNotFoundTeam        = errors.New("team not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrCannotParseToken    = errors.New("cannot parse token")
	ErrCannotGetUser       = errors.New("cannot get user")
	ErrCannotSetParam      = errors.New("cannot set param")
	ErrCannotSignToken     = errors.New("cannot sign token")
	ErrTokenExpired        = errors.New("token expired")
	ErrPRExists            = errors.New("pull request already exists")
	ErrCreatePR            = errors.New("cannot create pull request")
	ErrInactiveCreator     = errors.New("creator is inactive")
)
