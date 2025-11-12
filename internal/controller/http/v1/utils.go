package httpapi

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrors"

	"github.com/labstack/echo/v4"
)

func newErrReasonJSON(c echo.Context, httpCode int, msgCode httperrors.ErrorCode, msg string) error {
	return c.JSON(httpCode, httpdto.ErrorResponse{Error: httpdto.APIError{Code: msgCode, Message: msg}})
}
