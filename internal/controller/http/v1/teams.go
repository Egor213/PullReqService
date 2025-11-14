package httpapi

import (
	"errors"
	"net/http"

	"app/internal/controller/http/v1/httpdto"
	"app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	httpmappers "app/internal/controller/http/v1/mappers"
	mw "app/internal/controller/http/v1/midlleware"
	"app/internal/service"
	"app/internal/service/serverrs"

	"github.com/labstack/echo/v4"
)

type TeamsRoutes struct {
	teamsService service.Teams
}

func newTeamsRoutes(g *echo.Group, teamsServ service.Teams, m *mw.Auth) {
	r := &TeamsRoutes{
		teamsService: teamsServ,
	}

	g.POST("/add", r.addTeam)
	g.GET("/get", r.getTeam, m.UserIdentity)
}

func (r *TeamsRoutes) addTeam(c echo.Context) error {
	var input httpdto.AddTeamInput
	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
		return err
	}

	team, err := r.teamsService.CreateOrUpdateTeam(c.Request().Context(), httpmappers.ToCrOrUpTeamInput(input))
	if err != nil {
		if errors.Is(err, serverrs.ErrTeamWithUsersExists) {
			ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeTeamExists, err.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, httperrs.ErrCodeInternalServer, httperrs.ErrInternalServer.Error())
		return err
	}

	output := httpmappers.ToAddTeamOutput(team)
	return c.JSON(http.StatusCreated, output)
}

func (r *TeamsRoutes) getTeam(c echo.Context) error {
	var input httpdto.GetTeamInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, httperrs.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, httperrs.ErrCodeInvalidParams, err.Error())
		return err
	}

	team, err := r.teamsService.GetTeam(c.Request().Context(), input.TeamName)
	if err != nil {
		if errors.Is(err, serverrs.ErrNotFoundTeam) {
			ut.NewErrReasonJSON(c, http.StatusNotFound, httperrs.ErrCodeNotFound, httperrs.ErrNotFound.Error())
			return err

		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, httperrs.ErrCodeInternalServer, httperrs.ErrInternalServer.Error())
		return err

	}

	return c.JSON(http.StatusOK, team)
}
