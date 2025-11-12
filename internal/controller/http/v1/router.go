package httpapi

import (
	"app/internal/service"
	errorsutils "app/pkg/errors"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: setLogsFile()}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })

}

func setLogsFile() *os.File {
	logPath := filepath.Join("logs", "logfile.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(errorsutils.WrapPathErr(err))
	}
	return file
}
