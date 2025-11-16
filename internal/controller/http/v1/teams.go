package httpapi

import (
	"app/internal/service"
	"app/internal/usecase"
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

type TeamsRoutes struct {
	teamsService service.Teams
	teamsPRUC    usecase.TeamsPRUseCase
}

func newTeamsRoutes(g *echo.Group, teamsServ service.Teams, uc usecase.TeamsPRUseCase, m *mw.Auth) {
	r := &TeamsRoutes{
		teamsService: teamsServ,
		teamsPRUC:    uc,
	}

	g.POST("/add", r.addTeam)
	g.GET("/get", r.getTeam, m.UserIdentity, m.CheckRole(e.RoleUser))
	g.POST("/deactivate", r.deactivateTeam, m.UserIdentity, m.CheckRole(e.RoleAdmin))
}

func (r *TeamsRoutes) addTeam(c echo.Context) error {
	var input hd.AddTeamInput
	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	team, err := r.teamsService.CreateOrUpdateTeam(c.Request().Context(), hmap.ToCrOrUpTeamInput(input))
	if err != nil {
		if errors.Is(err, se.ErrTeamWithUsersExists) {
			return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeTeamExists, err.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}

	output := hmap.ToAddTeamOutput(team)
	return c.JSON(http.StatusCreated, output)
}

func (r *TeamsRoutes) getTeam(c echo.Context) error {
	var input hd.GetTeamInput

	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	team, err := r.teamsService.GetTeam(c.Request().Context(), input.TeamName)
	if err != nil {
		if errors.Is(err, se.ErrNotFoundTeam) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}
	output := hmap.ToGetTeamOutput(team)
	return c.JSON(http.StatusOK, output)
}

func (r *TeamsRoutes) deactivateTeam(c echo.Context) error {
	var input hd.DeactivateTeamInput

	if err := c.Bind(&input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	err := r.teamsPRUC.DeactivateTeamUsers(c.Request().Context(), input.TeamName)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrMergedPR):
			return ut.NewErrReasonJSON(c, http.StatusConflict, he.ErrCodePRMerged, he.ErrPRMerged.Error())
		case errors.Is(err, se.ErrNotFoundUser), errors.Is(err, se.ErrNotFoundPR):
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		default:
			return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		}
	}

	return c.JSON(http.StatusOK, hd.DeactivateTeamOutput(input))
}
