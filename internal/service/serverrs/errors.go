package serverrs

import "errors"

var (
	ErrTeamWithUsersExists = errors.New("team with such users already exists")
	ErrNotFoundTeam        = errors.New("team not found")
)
