package httpapi

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	"app/internal/service"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}

	g.POST("/login", r.login)
}

func (r *authRoutes) login(c echo.Context) error {
	var input httpdto.LoginInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), servdto.GenTokenInput{
		UserID: input.UserID,
		Role:   input.Role,
	})
	if err != nil {
		if errors.Is(err, serverrs.ErrUserNotFound) {
			ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeNotFound, httperrs.ErrNotFound.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, httperrs.ErrCodeInternalServer, httperrs.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusOK, httpdto.LoginOutput{
		AccessToken: token,
	})
}
