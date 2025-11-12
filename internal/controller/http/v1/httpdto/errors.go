package httpdto

import "app/internal/controller/http/v1/httperrors"

type APIError struct {
	Code    httperrors.ErrorCode `json:"code"`
	Message string               `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}
