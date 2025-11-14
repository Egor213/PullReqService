package httpapi

import (
	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	ut "app/internal/controller/http/v1/utils"
	"app/internal/service"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
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
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), sd.GenTokenInput{
		UserID: input.UserID,
		Role:   input.Role,
	})
	if err != nil {
		if errors.Is(err, se.ErrNotFoundUser) {
			return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}

	return c.JSON(http.StatusOK, hd.LoginOutput{
		AccessToken: token,
	})
}
