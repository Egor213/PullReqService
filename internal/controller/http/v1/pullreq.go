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

	g.GET("/get", r.getPR)
	g.POST("/create", r.createPR)
}

// TODO: если неактивный пользователь захочет создать PR?
func (r *PullReqRoutes) createPR(c echo.Context) error {
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
			ut.NewErrReasonJSON(c, http.StatusConflict, httperrs.ErrCodePRExists, err.Error())
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

func (r *PullReqRoutes) getPR(c echo.Context) error {
	var input httpdto.GetPRInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
		return err
	}

	pr, err := r.prService.GetPR(c.Request().Context(), input.PRID)
	if err != nil {
		if errors.Is(err, serverrs.ErrNotFoundPR) {
			ut.NewErrReasonJSON(c, http.StatusNotFound, httperrs.ErrCodeNotFound, httperrs.ErrNotFound.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, httperrs.ErrCodeInternalServer, httperrs.ErrInternalServer.Error())
		return err

	}

	return c.JSON(http.StatusOK, httpdto.GetPROutput{
		NeedMoreReviewers: &pr.NeedMoreReviewers,
		PullRequestDTO: httpdto.PullRequestDTO{
			PullReqID: pr.PullReqID,
			NamePR:    pr.NamePR,
			AuthorID:  pr.AuthorID,
			Status:    string(pr.Status),
			Reviewers: pr.Reviewers,
		},
	})
}
