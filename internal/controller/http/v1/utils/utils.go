package httputils

import (
	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

func NewErrReasonJSON(c echo.Context, httpCode int, msgCode he.ErrorCode, msg string) error {
	err := c.JSON(httpCode, hd.ErrorOutput{
		Error: hd.APIError{
			Code:    msgCode,
			Message: msg,
		},
	})
	if err != nil {
		return err
	}
	return fmt.Errorf("%s: %s", msgCode, msg)
}
