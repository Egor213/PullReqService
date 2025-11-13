package httperrs

import "errors"

type ErrorCode string

const (
	ErrCodeTeamExists     ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists       ErrorCode = "PR_EXISTS"
	ErrCodePRMerged       ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned    ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate    ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeInvalidParams  ErrorCode = "INVALID_REQUEST_PARAMETERS"
	ErrCodeInternalServer ErrorCode = "INTERVAL_SERVER_ERROR"
)

var (
	ErrInvalidParams  = errors.New("invalid request parameters")
	ErrNotFound       = errors.New("resource not found")
	ErrInternalServer = errors.New("internal server error")
)
