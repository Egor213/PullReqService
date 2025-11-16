package serverrs

import (
	"errors"
)

var (
	ErrNotFoundTeam      = errors.New("team not found")
	ErrNotFoundUser      = errors.New("user not found")
	ErrNotFoundReviewers = errors.New("reviewers not found")
	ErrNotFoundPR        = errors.New("pull request not found")

	ErrCannotParseToken         = errors.New("cannot parse token")
	ErrCannotGetUser            = errors.New("cannot get user")
	ErrCannotSetParam           = errors.New("cannot set param")
	ErrCannotSignToken          = errors.New("cannot sign token")
	ErrCannotGetPR              = errors.New("cannot get pull request")
	ErrCannotCreatePR           = errors.New("cannot create pull request")
	ErrCannotChangeReviewer     = errors.New("cannot change reviewer")
	ErrCannotChangeSetNeedMoRev = errors.New("cannot set more reviewers")
	ErrCannotAssignReviewers    = errors.New("cannot assign reviewers")
	ErrCannotGetTeam            = errors.New("cannot get team")
	ErrCannotCreateTeam         = errors.New("cannot create team")
	ErrCannotUpsetUsers         = errors.New("cannot update or create users")
	ErrCannotDelUsersFromTeam   = errors.New("cannot delete users from team")
	ErrCannotGetReviewerStats   = errors.New("cannot get reviewer stats")
	ErrCannotGetPRStats         = errors.New("cannot get pr stats")
	ErrCannotDelReviewer        = errors.New("cannot delete reviewer")
	ErrCannotDeactivateTeam     = errors.New("cannot deactivate users team")

	ErrPRExists            = errors.New("pull request already exists")
	ErrTeamWithUsersExists = errors.New("team with such users already exists")
	ErrTeamWExists         = errors.New("team already exists")

	ErrInactiveCreator      = errors.New("creator is inactive")
	ErrTokenExpired         = errors.New("token expired")
	ErrReviewerNotAssigned  = errors.New("specified reviewer is not assigned to this pull request")
	ErrNoAvailableReviewers = errors.New("no available reviewers to reassign")
	ErrMergedPR             = errors.New("pull request has been merged")
)
