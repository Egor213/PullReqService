package serverrs

import (
	"errors"
)

// TODO: добавь нормальные ошибки в методах сервисов, а то че там просто err надо через эти ошибки
var (
	ErrNotFoundTeam      = errors.New("team not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrNotFoundReviewers = errors.New("reviewers not found")
	ErrNotFoundPR        = errors.New("pull request not found")

	ErrCannotParseToken = errors.New("cannot parse token")
	ErrCannotGetUser    = errors.New("cannot get user")
	ErrCannotSetParam   = errors.New("cannot set param")
	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrCannotGetPR      = errors.New("cannot get pull request")
	ErrCreatePR         = errors.New("cannot create pull request")

	ErrPRExists            = errors.New("pull request already exists")
	ErrTeamWithUsersExists = errors.New("team with such users already exists")

	ErrInactiveCreator      = errors.New("creator is inactive")
	ErrTokenExpired         = errors.New("token expired")
	ErrReviewerNotAssigned  = errors.New("specified reviewer is not assigned to this pull request")
	ErrNoAvailableReviewers = errors.New("no available reviewers to reassign")
	ErrMergedPR             = errors.New("pull request has been merged")
)
