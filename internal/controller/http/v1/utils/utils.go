package httputils

import (
	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"

	"github.com/labstack/echo/v4"
)

func NewErrReasonJSON(c echo.Context, httpCode int, msgCode he.ErrorCode, msg string) error {
	return c.JSON(httpCode, hd.ErrorOutput{Error: hd.APIError{Code: msgCode, Message: msg}})
}
