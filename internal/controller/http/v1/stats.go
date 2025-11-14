package httpapi

import (
	"app/internal/service"
	"net/http"

	hd "app/internal/controller/http/v1/dto"
	he "app/internal/controller/http/v1/errors"
	hmap "app/internal/controller/http/v1/mappers"
	mw "app/internal/controller/http/v1/midlleware"
	ut "app/internal/controller/http/v1/utils"

	"github.com/labstack/echo/v4"
)

type StatsRoutes struct {
	statsService service.Stats
}

func newStatsRoutes(g *echo.Group, statsServ service.Stats, m *mw.Auth) {
	r := &StatsRoutes{
		statsService: statsServ,
	}

	g.GET("", r.getStats, m.UserIdentity)
}

func (r *StatsRoutes) getStats(c echo.Context) error {
	stats, err := r.statsService.GetStats(c.Request().Context())
	if err != nil {
		return ut.NewErrReasonJSON(c, http.StatusInternalServerError, he.ErrCodeForbidden, he.ErrInternalServer.Error())
	}

	return c.JSON(http.StatusOK, hd.GetStatsOutput{
		ByPRs:   hmap.ToPRStats(stats.ByPRs),
		ByUsers: hmap.ToReviewerStats(stats.ByUsers),
	})
}
