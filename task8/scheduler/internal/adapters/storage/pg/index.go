package pg

import (
	"context"
	// driver for migration
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	connectionWaitTimout = 30 * time.Second
	migrationWaitTimout  = 30 * time.Second
)

// PostgresDb - is type for postgres-db
type PostgresDb struct {
	db *sqlx.DB
}

// NewPostgresDb - creates new PostgresDb instance
func NewPostgresDb(dsn string) (*PostgresDb, error) {
	var err error

	res := &PostgresDb{}

	connectionContext, _ := context.WithTimeout(context.Background(), connectionWaitTimout)
	res.db, err = res.connectionWait(dsn, connectionContext)
	if err != nil {
		return nil, err
	}

	migrationContext, _ := context.WithTimeout(context.Background(), migrationWaitTimout)
	err = res.migrationWait(migrationContext)
	if err != nil {
		return nil, err
	}

	res.db.SetMaxOpenConns(10)
	res.db.SetMaxIdleConns(5)
	res.db.SetConnMaxLifetime(10 * time.Minute)

	return res, nil
}

func (pdb *PostgresDb) connectionWait(dsn string, ctx context.Context) (*sqlx.DB, error) {
	res, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	for {
		err = res.PingContext(ctx)
		if err == nil || err == ctx.Err() {
			break
		}
		time.Sleep(time.Second)
	}

	return res, err
}

func (pdb *PostgresDb) migrationWait(ctx context.Context) error {
	var err error
	var cnt uint32

	for {
		err := pdb.db.GetContext(ctx, &cnt, `select count(*) from event`)
		if err == nil || err == ctx.Err() {
			break
		}
		time.Sleep(time.Second)
	}

	return err
}
