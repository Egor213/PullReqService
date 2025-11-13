package httputils

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"

	"github.com/labstack/echo/v4"
)

func NewErrReasonJSON(c echo.Context, httpCode int, msgCode httperrs.ErrorCode, msg string) error {
	return c.JSON(httpCode, httpdto.ErrorOutput{Error: httpdto.APIError{Code: msgCode, Message: msg}})
}
