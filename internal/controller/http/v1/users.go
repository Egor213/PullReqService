package httpapi

import (
	"errors"
	"net/http"

	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"
	"app/internal/service"
	"app/internal/service/serverrs"

	"github.com/labstack/echo/v4"
)

type UsersRoutes struct {
	usersService service.Users
}

func newUsersRoutes(g *echo.Group, usersServ service.Users) {
	r := &UsersRoutes{
		usersService: usersServ,
	}

	g.POST("/setIsActive", r.setIsActive)
}

func (r *UsersRoutes) setIsActive(c echo.Context) error {
	var input httpdto.SetIsActiveInput
	if err := c.Bind(&input); err != nil {
		newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
		return err
	}

	user, err := r.usersService.SetIsActive(c.Request().Context(), input.UserID, input.IsActive)
	if err != nil {
		if errors.Is(err, serverrs.ErrUserNotFound) {
			newErrReasonJSON(c, http.StatusNotFound, httperrs.ErrCodeNotFound, httperrs.ErrNotFound.Error())
			return err
		}
		newErrReasonJSON(c, http.StatusInternalServerError, httperrs.ErrCodeInternalServer, httperrs.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusOK, httpdto.SetIsActiveOutput{User: user})
}
