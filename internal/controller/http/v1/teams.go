package httpapi

import (
	"errors"
	"net/http"

	hd "app/internal/controller/http/v1/httpdto"
	he "app/internal/controller/http/v1/httperrs"
	ut "app/internal/controller/http/v1/httputils"
	httpmappers "app/internal/controller/http/v1/mappers"
	mw "app/internal/controller/http/v1/midlleware"
	e "app/internal/entity"
	"app/internal/service"
	se "app/internal/service/serverrs"

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
	g.GET("/get", r.getTeam, m.UserIdentity, m.CheckRole(e.RoleUser))
}

func (r *TeamsRoutes) addTeam(c echo.Context) error {
	var input hd.AddTeamInput
	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	team, err := r.teamsService.CreateOrUpdateTeam(c.Request().Context(), httpmappers.ToCrOrUpTeamInput(input))
	if err != nil {
		if errors.Is(err, se.ErrTeamWithUsersExists) {
			ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeTeamExists, err.Error())
			return err
		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err
	}

	output := httpmappers.ToAddTeamOutput(team)
	return c.JSON(http.StatusCreated, output)
}

func (r *TeamsRoutes) getTeam(c echo.Context) error {
	var input hd.GetTeamInput

	if err := c.Bind(&input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
		return err
	}

	team, err := r.teamsService.GetTeam(c.Request().Context(), input.TeamName)
	if err != nil {
		if errors.Is(err, se.ErrNotFoundTeam) {
			ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
			return err

		}
		ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
		return err

	}
	output := httpmappers.ToGetTeamOutput(team)
	return c.JSON(http.StatusOK, output)
}
