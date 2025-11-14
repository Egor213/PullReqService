package httpapi

import (
	hd "app/internal/controller/http/v1/httpdto"
	he "app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	mw "app/internal/controller/http/v1/midlleware"
	"app/internal/service"
	sd "app/internal/service/servdto"
	se "app/internal/service/serverrs"
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
	g.POST("/reassign", r.reassignReviewer)
}

// TODO: если неактивный пользователь захочет создать PR?
func (r *PullReqRoutes) createPR(c echo.Context) error {
	var input hd.CreatePRInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	pr, err := r.prService.CreatePR(c.Request().Context(), sd.CreatePRInput{
		PullReqID: input.PullReqID,
		NamePR:    input.NamePR,
		AuthorID:  input.AuthorID,
	})

	if err != nil {
		if errors.Is(err, se.ErrPRExists) {
			ut.NewErrReasonJSON(c, http.StatusConflict, he.ErrCodePRExists, err.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusCreated, hd.CreatePROutput{
		PullReq: hd.PullRequestDTO{
			PullReqID: pr.PullReqID,
			NamePR:    pr.NamePR,
			AuthorID:  pr.AuthorID,
			Status:    string(pr.Status),
			Reviewers: pr.Reviewers,
		},
	})
}

func (r *PullReqRoutes) getPR(c echo.Context) error {
	var input hd.GetPRInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	pr, err := r.prService.GetPR(c.Request().Context(), input.PRID)
	if err != nil {
		if errors.Is(err, se.ErrNotFoundPR) {
			ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err

	}

	return c.JSON(http.StatusOK, hd.GetPROutput{
		NeedMoreReviewers: &pr.NeedMoreReviewers,
		PullRequestDTO: hd.PullRequestDTO{
			PullReqID: pr.PullReqID,
			NamePR:    pr.NamePR,
			AuthorID:  pr.AuthorID,
			Status:    string(pr.Status),
			Reviewers: pr.Reviewers,
		},
	})
}

func (r *PullReqRoutes) reassignReviewer(c echo.Context) error {
	var input hd.ReassignReviewerInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	out, err := r.prService.ReassignReviewer(c.Request().Context(), sd.ReassignReviewerInput{
		PullReqID: input.PullReqID,
		RevID:     input.OldReviewer,
	})

	if err != nil {

		if errors.Is(err, se.ErrReviewerNotAssigned) || errors.Is(err, se.ErrNoAvailableReviewers) {
			ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
			return err
		} else if errors.Is(err, se.ErrMergedPR) {
			ut.NewErrReasonJSON(c, http.StatusConflict, he.ErrCodePRMerged, he.ErrPRMerged.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err
	}

	return c.JSON(http.StatusOK, hd.ReassignReviewerOutput{
		PullReq: hd.PullRequestDTO{
			PullReqID: out.PullReq.PullReqID,
			NamePR:    out.PullReq.NamePR,
			AuthorID:  out.PullReq.AuthorID,
			Status:    string(out.PullReq.Status),
			Reviewers: out.PullReq.Reviewers,
		},
		NewReviewer: out.NewRevID,
	})
}
