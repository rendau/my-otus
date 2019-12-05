package pg

import (
	"errors"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

// PostgresDb - is type for memory-db
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

// MigrationDo - applies or reverts migrations
func MigrationDo(dsn, migrationPath, cmd string) error {
	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		return err
	}

	switch cmd {
	case "up":
		err = m.Up()
		break
	case "down":
		err = m.Down()
		break
	default:
		return errors.New("bad command")
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
