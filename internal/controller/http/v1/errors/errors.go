package httperrs

import (
	"errors"
)

type ErrorCode string

const (
	ErrCodeTeamExists      ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists        ErrorCode = "PR_EXISTS"
	ErrCodePRMerged        ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned     ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate     ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeInvalidParams   ErrorCode = "INVALID_REQUEST_PARAMETERS"
	ErrCodeInternalServer  ErrorCode = "INTERVAL_SERVER_ERROR"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeInactiveCreator ErrorCode = "INACTIVE_CREATOR"
)

var (
	ErrInvalidParams     = errors.New("invalid request parameters")
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInternalServer    = errors.New("internal server error")
	ErrNoRights          = errors.New("no rights")
	ErrPRMerged          = errors.New("cannot reassign on merged PR")
	ErrPRAlreadyExists   = errors.New("PR id already exists")
	ErrReviewerNotAssign = errors.New("reviewer is not assigned to this PR")
	ErrNoActiveCandidate = errors.New("no active replacement candidate in team")
)
