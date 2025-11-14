package httpapi

import (
	"errors"
	"net/http"

	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	hmap "app/internal/controller/http/v1/mappers"
	mw "app/internal/controller/http/v1/midlleware"
	ut "app/internal/controller/http/v1/utils"
	e "app/internal/entity"
	"app/internal/service"
	sd "app/internal/service/servdto"
	se "app/internal/service/serverrs"
	errutils "app/pkg/errors"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type UsersRoutes struct {
	uService  service.Users
	prService service.PullReq
}

func newUsersRoutes(g *echo.Group, uServ service.Users, prServ service.PullReq, m *mw.Auth) {
	r := &UsersRoutes{
		uService:  uServ,
		prService: prServ,
	}

	g.GET("/getReview", r.getReview, m.UserIdentity, m.CheckRole(e.RoleUser))
	g.POST("/setIsActive", r.setIsActive, m.UserIdentity, m.CheckRole(e.RoleAdmin))
}

func (r *UsersRoutes) setIsActive(c echo.Context) error {
	var input hd.SetIsActiveInput
	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	user, err := r.uService.SetIsActive(c.Request().Context(), sd.SetIsActiveInput{
		UserID:   input.UserID,
		IsActive: input.IsActive,
	})
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		if errors.Is(err, se.ErrUserNotFound) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}

	return c.JSON(http.StatusOK, hd.SetIsActiveOutput{
		User: user,
	})
}

func (r *UsersRoutes) getReview(c echo.Context) error {
	var input hd.GetReviewInput

	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	prs, err := r.prService.GetPRsByReviewer(c.Request().Context(), input.UserID)
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		if errors.Is(err, se.ErrUserNotFound) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}
	output := hmap.ToGetReviewOutput(input.UserID, prs)
	return c.JSON(http.StatusOK, output)
}
