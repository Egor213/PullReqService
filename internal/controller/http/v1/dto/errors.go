package httpdto

import httperrs "app/internal/controller/http/v1/errors"

type APIError struct {
	Code    httperrs.ErrorCode `json:"code"`
	Message string             `json:"message"`
}

type ErrorOutput struct {
	Error APIError `json:"error"`
}
