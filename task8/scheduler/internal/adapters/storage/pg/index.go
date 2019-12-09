package pg

import (
	// driver for migration
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

// PostgresDb - is type for postgres-db
type PostgresDb struct {
	db *sqlx.DB
}

// NewPostgresDb - creates new PostgresDb instance
func NewPostgresDb(dsn string) (*PostgresDb, error) {
	var err error

	res := &PostgresDb{}

	res.db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	res.db.SetMaxOpenConns(10)
	res.db.SetMaxIdleConns(5)
	res.db.SetConnMaxLifetime(10 * time.Minute)

	return res, nil
}
