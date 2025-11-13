package httpapi

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	mw "app/internal/controller/http/v1/midlleware"
	"app/internal/service"
	"app/internal/service/servdto"
	"app/internal/service/serverrs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PullReqRoutes struct {
	prService service.PullReq
}

func newPullReqRoutes(g *echo.Group, prService service.PullReq, m *mw.Auth) {
	r := &PullReqRoutes{
		prService: prService,
	}

	g.POST("/create", r.CreatePR)
}

func (r *PullReqRoutes) CreatePR(c echo.Context) error {
	var input httpdto.CreatePRInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
		return err
	}

	pr, err := r.prService.CreatePR(c.Request().Context(), servdto.CreatePRInput{
		PullReqID: input.PullReqID,
		NamePR:    input.NamePR,
		AuthorID:  input.AuthorID,
	})

	if err != nil {
		if errors.Is(err, serverrs.ErrPRExists) {
			ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodePRExists, err.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, httperrs.ErrCodeInternalServer, httperrs.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusCreated, httpdto.CreatePROutput{
		PullReq: httpdto.PullRequestDTO{
			PullReqID: pr.PullReqID,
			NamePR:    pr.NamePR,
			AuthorID:  pr.AuthorID,
			Status:    string(pr.Status),
			Reviewers: pr.Reviewers,
		},
	})
}
