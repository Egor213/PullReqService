package mw

import (
	"app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	e "app/internal/entity"
	"app/internal/service"
	errorsutils "app/pkg/errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const (
	userIdCtx = "UserId"
	roleKey   = "Role"
)

type Auth struct {
	authService service.Auth
}

func NewAuth(authService service.Auth) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (h *Auth) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			log.Errorf("bearerToken: %v", errorsutils.WrapPathErr(httperrs.ErrInvalidAuthHeader))
			ut.NewErrReasonJSON(c, http.StatusUnauthorized, httperrs.ErrCodeInvalidHeader, httperrs.ErrInvalidAuthHeader.Error())
			return nil
		}

		claims, err := h.authService.ParseToken(token)
		if err != nil {
			log.Error(errorsutils.WrapPathErr(err).Error())
			ut.NewErrReasonJSON(c, http.StatusUnauthorized, httperrs.ErrCodeInvalidToken, httperrs.ErrCannotParseToken.Error())
			return err
		}

		c.Set(userIdCtx, claims.UserID)
		c.Set(roleKey, claims.Role)

		return next(c)

	}
}

func (m *Auth) CheckRole(required e.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(roleKey).(e.Role)
			if !ok || role != required {
				return echo.NewHTTPError(http.StatusForbidden, httperrs.ErrNoRights.Error())
			}

			return next(c)
		}
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
