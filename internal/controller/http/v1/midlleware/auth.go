package mw

import (
	he "app/internal/controller/http/v1/errors"
	ut "app/internal/controller/http/v1/utils"
	e "app/internal/entity"
	"app/internal/service"
	errutils "app/pkg/errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const (
	userIdCtx = "UserId"
	roleKey   = "Role"
)

var rolePriority = map[e.Role]int{
	e.RoleUser:  1,
	e.RoleAdmin: 2,
}

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
			log.Error(errutils.WrapPathErr(fmt.Errorf("invalid auth header")))
			return ut.NewErrReasonJSON(c, http.StatusUnauthorized, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}

		claims, err := h.authService.ParseToken(token)
		if err != nil {
			return ut.NewErrReasonJSON(c, http.StatusUnauthorized, he.ErrCodeNotFound, he.ErrNotFound.Error())
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
			if !ok || rolePriority[role] < rolePriority[required] {
				return ut.NewErrReasonJSON(c, http.StatusForbidden, he.ErrCodeForbidden, he.ErrNoRights.Error())
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
