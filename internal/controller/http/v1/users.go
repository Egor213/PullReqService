package httpapi

import (
	"errors"
	"net/http"

	hd "app/internal/controller/http/v1/httpdto"
	he "app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	mw "app/internal/controller/http/v1/midlleware"
	e "app/internal/entity"
	"app/internal/service"
	sd "app/internal/service/servdto"
	se "app/internal/service/serverrs"

	"github.com/labstack/echo/v4"
)

type UsersRoutes struct {
	usersService service.Users
}

func newUsersRoutes(g *echo.Group, usersServ service.Users, m *mw.Auth) {
	r := &UsersRoutes{
		usersService: usersServ,
	}

	g.POST("/setIsActive", r.setIsActive, m.UserIdentity, m.CheckRole(e.RoleAdmin))
}

func (r *UsersRoutes) setIsActive(c echo.Context) error {
	var input hd.SetIsActiveInput
	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	user, err := r.usersService.SetIsActive(c.Request().Context(), sd.SetIsActiveInput{
		UserID:   input.UserID,
		IsActive: input.IsActive,
	})
	if err != nil {
		if errors.Is(err, se.ErrUserNotFound) {
			ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusOK, hd.SetIsActiveOutput{User: user})
}
