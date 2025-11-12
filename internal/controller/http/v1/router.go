package httpapi

import (
	"app/internal/service"
	errorsutils "app/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	logFile := setLogsFile()
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: multiWriter,
	}))

	handler.Use(middleware.Recover())

	api := handler.Group("/api")
	{
		api.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

		teams := api.Group("/teams")
		newTeamsRoutes(teams, services.Teams)

	}

}

func setLogsFile() *os.File {
	logPath := filepath.Join("logs", "logfile.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(errorsutils.WrapPathErr(err))
	}
	return file
}
