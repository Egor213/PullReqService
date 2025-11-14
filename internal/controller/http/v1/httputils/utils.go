package httputils

import (
	hd "app/internal/controller/http/v1/httpdto"
	he "app/internal/controller/http/v1/httperrs"

	"github.com/labstack/echo/v4"
)

func NewErrReasonJSON(c echo.Context, httpCode int, msgCode he.ErrorCode, msg string) error {
	return c.JSON(httpCode, hd.ErrorOutput{Error: hd.APIError{Code: msgCode, Message: msg}})
}
