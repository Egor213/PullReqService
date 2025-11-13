package httpdto

import "app/internal/controller/http/v1/httperrs"

type APIError struct {
	Code    httperrs.ErrorCode `json:"code"`
	Message string             `json:"message"`
}

type ErrorOutput struct {
	Error APIError `json:"error"`
}
