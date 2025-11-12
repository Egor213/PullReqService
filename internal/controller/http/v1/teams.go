package httpapi

import (
	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrors"
	"app/internal/service"
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
		return newErrReasonJSON(c, http.StatusBadRequest, httperrors.ErrCodeInvalidParams, httperrors.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, httperrors.ErrCodeInvalidParams, err.Error())
	}

	return c.JSON(http.StatusOK, httpdto.AddTeamOutput{
		Team: input,
	})
}
