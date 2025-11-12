package httpapi

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"
	"app/internal/controller/http/v1/mappers"
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
}

func (r *TeamsRoutes) addTeam(c echo.Context) error {
	var input httpdto.AddTeamInput
	if err := c.Bind(&input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
	}

	team, err := r.teamsService.CreateOrUpdateTeam(c.Request().Context(), mappers.ToEntityTeam(input))

	if err != nil {
		if errors.Is(err, serverrs.ErrTeamWithUsersExists) {
			return newErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeTeamExists, err.Error())
		}
		return err
	}

	output := mappers.ToAddTeamOutput(team)
	return c.JSON(http.StatusOK, output)
}
