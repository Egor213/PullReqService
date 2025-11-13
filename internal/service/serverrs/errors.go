package serverrs

import (
	"errors"
)

var (
	ErrTeamWithUsersExists = errors.New("team with such users already exists")
	ErrNotFoundTeam        = errors.New("team not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrCannotParseToken    = errors.New("cannot parse token")
	ErrCannotGetUser       = errors.New("cannot get user")
	ErrCannotSetParam      = errors.New("cannot set param")
	ErrCannotSignToken     = errors.New("cannot sign token")
)
