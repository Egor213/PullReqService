package app

import (
	"app/internal/config"
	httpapi "app/internal/controller/http/v1"
	"app/internal/repo"
	"app/internal/service"
	errutils "app/pkg/errors"
	"app/pkg/httpserver"
	"app/pkg/logger"
	"app/pkg/postgres"
	"app/pkg/validator"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func Run() {
	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatal(errutils.WrapPathErr(err))
	}

	// Logger
	logger.SetupLogger(cfg.Log.Level)
	log.Info("Logger has been set up")

	// Migrations
	Migrate(cfg.PG.URL)

	// DB connecting
	log.Info("Connecting to DB...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(errutils.WrapPathErr(err))
	}
	defer pg.Close()
	log.Info("Connected to DB")

	// Repos
	repositories := repo.NewRepositories(pg)

	// Services
	deps := service.ServicesDependencies{
		Repos: repositories,
	}
	services := service.NewServices(deps)

	// Echo handler
	log.Info("Initializing handlers and routes")
	handler := echo.New()

	handler.Validator = validator.NewCustomValidator()
	httpapi.ConfigureRouter(handler, services)

	// HTTP server
	log.Info("Starting http server")
	log.Debugf("Server port: %s", cfg.HTTP.Address)
	httpServer := httpserver.New(handler, httpserver.Address(cfg.HTTP.Address))

	// Waiting signal
	log.Info("Configuring graceful shutdown")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(errutils.WrapPathErr(err))
	}

	// Graceful shutdown
	log.Info("Shutting down")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
	}

}
