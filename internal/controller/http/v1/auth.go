package httpapi

import (
	hd "app/internal/controller/http/v1/httpdto"
	he "app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	"app/internal/service"
	sd "app/internal/service/servdto"
	se "app/internal/service/serverrs"
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
	var input hd.LoginInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), sd.GenTokenInput{
		UserID: input.UserID,
		Role:   input.Role,
	})
	if err != nil {
		if errors.Is(err, se.ErrUserNotFound) {
			ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeNotFound, he.ErrNotFound.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusOK, hd.LoginOutput{
		AccessToken: token,
	})
}
