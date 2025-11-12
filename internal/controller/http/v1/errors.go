package httpapi

import (
	"fmt"
)

var (
	ErrInvalidAuthHeader   = fmt.Errorf("invalid auth header")
	ErrCannotParseToken    = fmt.Errorf("cannot parse token")
	ErrInvalidBody         = fmt.Errorf("invalid request body")
	ErrInternalServerError = fmt.Errorf("internal server error")
)
