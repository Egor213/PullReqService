package httpapi

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"
	httpmappers "app/internal/controller/http/v1/mappers"
	"app/internal/service"
	"app/internal/service/serverrs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TeamsRoutes struct {
	teamsService service.Teams
}

func newTeamsRoutes(g *echo.Group, teamsServ service.Teams) {
	r := &TeamsRoutes{
		teamsService: teamsServ,
	}

	g.POST("/add", r.addTeam)
	g.GET("/get/", r.getTeam)
}

func (r *TeamsRoutes) addTeam(c echo.Context) error {
	var input httpdto.AddTeamInput
	if err := c.Bind(&input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
	}

	team, err := r.teamsService.CreateOrUpdateTeam(c.Request().Context(), httpmappers.ToEntityTeam(input))

	if err != nil {
		if errors.Is(err, serverrs.ErrTeamWithUsersExists) {
			return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeTeamExists, err.Error())
		}
		return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeTeamExists, httperrs.ErrInternalServer.Error())
	}

	output := httpmappers.ToAddTeamOutput(team)
	return c.JSON(http.StatusOK, output)
}

func (r *TeamsRoutes) getTeam(c echo.Context) error {
	var input httpdto.GetTeamInput

	if err := c.Bind(&input); err != nil {
		return newErrReasonJSON(
			c,
			http.StatusBadRequest,
			httperrs.ErrCodeInvalidParams,
			httperrs.ErrInvalidParams.Error(),
		)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(
			c,
			http.StatusBadRequest,
			httperrs.ErrCodeInvalidParams,
			err.Error(),
		)
	}

	team, err := r.teamsService.GetTeam(c.Request().Context(), input.TeamName)
	if err != nil {
		if errors.Is(err, serverrs.ErrNotFoundTeam) {
			return newErrReasonJSON(
				c,
				http.StatusNotFound,
				httperrs.ErrCodeNotFound,
				httperrs.ErrNotFound.Error(),
			)
		}
		return newErrReasonJSON(
			c,
			http.StatusBadRequest,
			httperrs.ErrCodeInternalServer,
			httperrs.ErrInternalServer.Error(),
		)
	}

	return c.JSON(http.StatusOK, team)
}
