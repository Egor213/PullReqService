package httpapi

import (
	"app/internal/service"
	"app/internal/usecase"
	ucd "app/internal/usecase/dto"
	"errors"
	"net/http"

	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	hmap "app/internal/controller/http/v1/mappers"
	mw "app/internal/controller/http/v1/midlleware"
	ut "app/internal/controller/http/v1/utils"
	e "app/internal/entity"

	se "app/internal/service/errors"

	"github.com/labstack/echo/v4"
)

type UsersRoutes struct {
	uService  service.Users
	prService service.PullReq
	usPRUC    usecase.UsersPRUseCase
}

func newUsersRoutes(g *echo.Group, uServ service.Users, prServ service.PullReq, uc usecase.UsersPRUseCase, m *mw.Auth) {
	r := &UsersRoutes{
		uService:  uServ,
		prService: prServ,
		usPRUC:    uc,
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

	user, err := r.usPRUC.SetIsActiveAndReassignPRs(c.Request().Context(), ucd.ActiveAndReassugnInput{
		UserID:   input.UserID,
		IsActive: input.IsActive,
	})

	if err != nil {
		if errors.Is(err, se.ErrNotFoundUser) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		} else if errors.Is(err, se.ErrMergedPR) {
			return ut.NewErrReasonJSON(c, http.StatusConflict, he.ErrCodePRMerged, he.ErrPRMerged.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}

	return c.JSON(http.StatusOK, hd.SetIsActiveOutput{
		User: hd.UserDTO{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
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
		if errors.Is(err, se.ErrNotFoundUser) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}
	output := hmap.ToGetReviewOutput(input.UserID, prs)
	return c.JSON(http.StatusOK, output)
}
