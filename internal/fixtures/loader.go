package fixtures

import (
	"database/sql"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Fixtures struct {
	loader *testfixtures.Loader
}

func NewFixtures(db *sql.DB, cfgPath, dialect string) (*Fixtures, error) {
	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dialect),
		testfixtures.Directory(cfgPath),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		return nil, err
	}

	return &Fixtures{loader: loader}, nil
}

func NewFixturesFromPool(pool *pgxpool.Pool, cfgPath, dialect string) (*Fixtures, error) {
	config := pool.Config().ConnConfig

	db := stdlib.OpenDB(*config)

	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dialect),
		testfixtures.Directory(cfgPath),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		return nil, err
	}

	return &Fixtures{loader: loader}, nil
}

func (f *Fixtures) Load() error {
	return f.loader.Load()
}
