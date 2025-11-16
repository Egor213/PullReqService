package main

import (
	"app/internal/config"
	"app/internal/fixtures"
	errutils "app/pkg/errors"
	"database/sql"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(errutils.WrapPathErr(err))
	}
	url := cfg.PG.URL + "?sslmode=disable"
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	fx, err := fixtures.NewFixtures(db, "./fixtures", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	if err := fx.Load(); err != nil {
		log.Fatal(err)
	}

	log.Println("Fixtures successfully loaded")
}
