package httpapi

import (
	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	mw "app/internal/controller/http/v1/midlleware"
	ut "app/internal/controller/http/v1/utils"
	e "app/internal/entity"
	"app/internal/service"
	sd "app/internal/service/dto"
	se "app/internal/service/errors"
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

	g.GET("/get", r.getPR, m.UserIdentity, m.CheckRole(e.RoleAdmin))
	g.POST("/merge", r.mergePR, m.UserIdentity, m.CheckRole(e.RoleAdmin))
	g.POST("/create", r.createPR, m.UserIdentity, m.CheckRole(e.RoleAdmin))
	g.POST("/reassign", r.reassignReviewer, m.UserIdentity, m.CheckRole(e.RoleAdmin))
}

func (r *PullReqRoutes) createPR(c echo.Context) error {
	var input hd.CreatePRInput

	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	pr, err := r.prService.CreatePR(c.Request().Context(), sd.CreatePRInput{
		PullReqID: input.PullReqID,
		NamePR:    input.NamePR,
		AuthorID:  input.AuthorID,
	})

	if err != nil {
		if errors.Is(err, se.ErrPRExists) {
			return ut.NewErrReasonJSON(c, http.StatusConflict, he.ErrCodePRExists, he.ErrPRAlreadyExists.Error())
		} else if errors.Is(err, se.ErrInactiveCreator) {
			return ut.NewErrReasonJSON(c, http.StatusForbidden, he.ErrCodeInactiveCreator, err.Error())
		} else if errors.Is(err, se.ErrNotFoundUserForPr) || errors.Is(err, se.ErrNotFoundTeam) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
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
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	pr, err := r.prService.GetPR(c.Request().Context(), input.PRID)
	if err != nil {
		if errors.Is(err, se.ErrNotFoundPR) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())

	}

	return c.JSON(http.StatusOK, hd.GetPROutput{
		NeedMoreReviewers: &pr.NeedMoreReviewers,
		MergedAt:          pr.MergedAt,
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
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	out, err := r.prService.ReassignReviewer(c.Request().Context(), sd.ReassignReviewerInput{
		PullReqID: input.PullReqID,
		RevID:     input.OldReviewer,
	})

	if err != nil {
		if errors.Is(err, se.ErrReviewerNotAssigned) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotAssigned, he.ErrReviewerNotAssign.Error())
		} else if errors.Is(err, se.ErrNoAvailableReviewers) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNoCandidate, he.ErrNoActiveCandidate.Error())
		} else if errors.Is(err, se.ErrMergedPR) {
			return ut.NewErrReasonJSON(c, http.StatusConflict, he.ErrCodePRMerged, he.ErrPRMerged.Error())
		} else if errors.Is(err, se.ErrNotFoundUser) || errors.Is(err, se.ErrNotFoundPR) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
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

func (r *PullReqRoutes) mergePR(c echo.Context) error {
	var input hd.MergePRInput

	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	pr, err := r.prService.MergePR(c.Request().Context(), input.PullReqID)
	if err != nil {
		if errors.Is(err, se.ErrNotFoundPR) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}
	return c.JSON(http.StatusOK, hd.MergePROutput{
		PullReq: hd.MergePRODTO{
			PullRequestDTO: hd.PullRequestDTO{
				PullReqID: pr.PullReqID,
				NamePR:    pr.NamePR,
				AuthorID:  pr.AuthorID,
				Status:    string(pr.Status),
				Reviewers: pr.Reviewers,
			},
			MergedAt: pr.MergedAt,
		},
	})
}
