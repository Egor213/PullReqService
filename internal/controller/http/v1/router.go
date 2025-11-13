package httpapi

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"app/internal/service"
	errorsutils "app/pkg/errors"

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
	handler.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	api := handler.Group("/api/v1")
	{
		newTeamsRoutes(api.Group("/team"), services.Teams)

	}
}

func setLogsFile() *os.File {
	logPath := filepath.Join("logs", "logfile.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		log.Fatal(errorsutils.WrapPathErr(err))
	}
	return file
}
