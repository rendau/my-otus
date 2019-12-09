package main

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/adapters/logger"
	"github.com/rendau/my-otus/task8/api/internal/adapters/storage/pg"
	"github.com/rendau/my-otus/task8/api/internal/domain/entities"
	"github.com/rendau/my-otus/task8/api/internal/domain/errors"
	"github.com/rendau/my-otus/task8/api/internal/domain/usecases"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

const testConfPath = "./test_conf.yml"

var ucs *usecases.Usecases

func TestMain(m *testing.M) {
	viper.SetConfigFile(testConfPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Fail to parse config")
	}

	lg, err := logger.NewLogger("", "", true, true)
	if err != nil {
		log.Fatalln("Fail to create logger")
	}
	defer lg.Sync()

	err = pg.MigrationDo(
		viper.GetString("pg_dsn"),
		viper.GetString("pg_migrations_path"),
		"down",
	)
	if err != nil {
		log.Fatalln("Fail to apply migrations, error:", err)
	}

	err = pg.MigrationDo(
		viper.GetString("pg_dsn"),
		viper.GetString("pg_migrations_path"),
		"up",
	)
	if err != nil {
		log.Fatalln("Fail to apply migrations, error:", err)
	}

	db, err := pg.NewPostgresDb(viper.GetString("pg_dsn"))
	if err != nil {
		log.Fatalln("Fail to create postgres-db, error:", err)
	}

	ucs = usecases.CreateUsecases(lg, db)

	// Start tests
	code := m.Run()

	err = pg.MigrationDo(
		viper.GetString("pg_dsn"),
		viper.GetString("pg_migrations_path"),
		"down",
	)
	if err != nil {
		log.Fatalln("Fail to apply migrations, error:", err)
	}

	os.Exit(code)
}

func TestList(t *testing.T) {
	var l []*entities.Event

	ctx := context.Background()

	var eCnt int // current count for comparison

	l, err := ucs.Event.List(ctx, nil)
	require.Nil(t, err)
	eCnt = len(l)

	e := &entities.Event{
		Owner:     "U1",
		Title:     "e1",
		Text:      "text",
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}
	err = ucs.Event.Create(ctx, e)
	require.Nil(t, err)

	e1, err := ucs.Event.Get(ctx, e.ID)
	require.Nil(t, err)
	require.NotNil(t, e1)
	require.Equal(t, e.Owner, e1.Owner)
	require.Equal(t, e.Title, e1.Title)
	require.Equal(t, e.Text, e1.Text)
	require.Equal(t, e.StartTime.UTC(), e1.StartTime.UTC())
	require.Equal(t, e.EndTime.UTC(), e1.EndTime.UTC())

	l, err = ucs.Event.List(ctx, nil)
	require.Nil(t, err)
	require.Equal(t, eCnt+1, len(l))

	st := time.Now().Add(3 * time.Hour)
	l, err = ucs.Event.List(ctx, &entities.EventListFilter{StartTimeGt: &st})
	require.Nil(t, err)
	require.Equal(t, eCnt, len(l))

	st = time.Now().Add(time.Minute)
	l, err = ucs.Event.List(ctx, &entities.EventListFilter{StartTimeLt: &st})
	require.Nil(t, err)
	require.Equal(t, eCnt, len(l))

	st = time.Now().Add(3 * time.Hour)
	l, err = ucs.Event.List(ctx, &entities.EventListFilter{EndTimeGt: &st})
	require.Nil(t, err)
	require.Equal(t, eCnt, len(l))

	st = time.Now().Add(time.Minute)
	l, err = ucs.Event.List(ctx, &entities.EventListFilter{EndTimeLt: &st})
	require.Nil(t, err)
	require.Equal(t, eCnt, len(l))

	st = time.Now().Add(90 * time.Minute)
	l, err = ucs.Event.List(ctx, &entities.EventListFilter{StartTimeLt: &st, EndTimeGt: &st})
	require.Nil(t, err)
	require.Equal(t, eCnt+1, len(l))

	err = ucs.Event.Delete(ctx, e.ID)
	require.Nil(t, err)

	l, err = ucs.Event.List(ctx, nil)
	require.Nil(t, err)
	require.Equal(t, 0, len(l))
}

func TestCreate(t *testing.T) {
	var err error

	ctx := context.Background()

	// ideal case
	e := &entities.Event{
		Owner:     "U1",
		Title:     "e1",
		Text:      "text",
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}
	err = ucs.Event.Create(ctx, e)
	require.Nil(t, err)

	err = ucs.Event.Delete(ctx, e.ID)
	require.Nil(t, err)

	type errCaseSt struct {
		owner string
		title string
		text  string
		st    time.Time
		et    time.Time
		err   error
	}
	errCases := []errCaseSt{
		errCaseSt{
			"",
			"e",
			"text",
			time.Now().Add(time.Hour),
			time.Now().Add(2 * time.Hour),
			errors.ErrOwnerRequired,
		},
		errCaseSt{
			"o",
			"",
			"text",
			time.Now().Add(time.Hour),
			time.Now().Add(2 * time.Hour),
			errors.ErrTitleRequired,
		},
		errCaseSt{
			"o",
			"e",
			"text",
			time.Now().Add(-time.Hour),
			time.Now().Add(2 * time.Hour),
			errors.ErrIncorrectStartDate,
		},
		errCaseSt{
			"o",
			"e",
			"text",
			time.Now().Add(2 * time.Hour),
			time.Now().Add(time.Hour),
			errors.ErrEndDateLTStartDate,
		},
		errCaseSt{
			"o",
			"e",
			"text",
			time.Now().Add(time.Minute),
			time.Now().Add(time.Hour),
			nil,
		},
		errCaseSt{
			"o",
			"e",
			"text",
			time.Now().Add(2 * time.Minute),
			time.Now().Add(5 * time.Minute),
			errors.ErrOverlaping,
		},
	}

	for _, c := range errCases {
		err = ucs.Event.Create(ctx, &entities.Event{
			Owner:     c.owner,
			Title:     c.title,
			Text:      c.text,
			StartTime: c.st,
			EndTime:   c.et,
		})
		if c.err != nil {
			require.NotNil(t, err)
			require.Equal(t, c.err, err)
		}
	}
}
