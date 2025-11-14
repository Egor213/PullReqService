package httpapi

import (
	"errors"
	"net/http"

	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	httpmappers "app/internal/controller/http/v1/mappers"
	mw "app/internal/controller/http/v1/midlleware"
	ut "app/internal/controller/http/v1/utils"
	e "app/internal/entity"
	"app/internal/service"
	se "app/internal/service/errors"
	errutils "app/pkg/errors"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, he.ErrInvalidParams.Error())
	}

	if err := c.Validate(input); err != nil {
		return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeInvalidParams, err.Error())
	}

	team, err := r.teamsService.CreateOrUpdateTeam(c.Request().Context(), httpmappers.ToCrOrUpTeamInput(input))
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		if errors.Is(err, se.ErrTeamWithUsersExists) {
			return ut.NewErrReasonJSON(c, http.StatusBadRequest, he.ErrCodeTeamExists, err.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}

	output := httpmappers.ToAddTeamOutput(team)
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
		log.Error(errutils.WrapPathErr(err))
		if errors.Is(err, se.ErrNotFoundTeam) {
			return ut.NewErrReasonJSON(c, http.StatusNotFound, he.ErrCodeNotFound, he.ErrNotFound.Error())
		}
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeInternalServer, he.ErrInternalServer.Error())
	}
	output := httpmappers.ToGetTeamOutput(team)
	return c.JSON(http.StatusOK, output)
}
