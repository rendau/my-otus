package main

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/adapters/storage/memdb"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/errors"
	"github.com/rendau/my-otus/task8/internal/domain/usecases"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	var err error
	var st time.Time
	var et time.Time
	var e *entities.Event
	var l []*entities.Event

	db, err := memdb.NewMemDb()
	require.Nil(t, err)

	uks := usecases.CreateUsecases(db)

	ctx := context.Background()

	st = time.Now().Add(time.Hour)
	et = time.Now().Add(2 * time.Hour)
	e, err = uks.Event.Create(ctx, "U1", "e1", "text", st, et)
	require.Nil(t, err)

	l, err = uks.Event.List(ctx, nil)
	require.Nil(t, err)
	require.Equal(t, 1, len(l))

	st = time.Now().Add(3 * time.Hour)
	l, err = uks.Event.List(ctx, &entities.EventListFilter{StartTimeGt: &st})
	require.Nil(t, err)
	require.Equal(t, 0, len(l))

	st = time.Now().Add(time.Minute)
	l, err = uks.Event.List(ctx, &entities.EventListFilter{StartTimeLt: &st})
	require.Nil(t, err)
	require.Equal(t, 0, len(l))

	st = time.Now().Add(3 * time.Hour)
	l, err = uks.Event.List(ctx, &entities.EventListFilter{EndTimeGt: &st})
	require.Nil(t, err)
	require.Equal(t, 0, len(l))

	st = time.Now().Add(time.Minute)
	l, err = uks.Event.List(ctx, &entities.EventListFilter{EndTimeLt: &st})
	require.Nil(t, err)
	require.Equal(t, 0, len(l))

	st = time.Now().Add(90 * time.Minute)
	l, err = uks.Event.List(ctx, &entities.EventListFilter{StartTimeLt: &st, EndTimeGt: &st})
	require.Nil(t, err)
	require.Equal(t, 1, len(l))

	err = uks.Event.Delete(ctx, e.ID)
	require.Nil(t, err)

	l, err = uks.Event.List(ctx, nil)
	require.Nil(t, err)
	require.Equal(t, 0, len(l))
}

func TestCreate(t *testing.T) {
	var err error
	var st time.Time
	var et time.Time
	var e *entities.Event

	db, err := memdb.NewMemDb()
	require.Nil(t, err)

	uks := usecases.CreateUsecases(db)

	ctx := context.Background()

	// ideal case
	st = time.Now().Add(time.Hour)
	et = time.Now().Add(2 * time.Hour)
	e, err = uks.Event.Create(ctx, "U1", "e1", "text", st, et)
	require.Nil(t, err)
	require.Equal(t, "U1", e.Owner)
	require.Equal(t, "e1", e.Title)
	require.Equal(t, "text", e.Text)
	require.Equal(t, st, e.StartTime)
	require.Equal(t, et, e.EndTime)

	err = uks.Event.Delete(ctx, e.ID)
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
		e, err = uks.Event.Create(ctx, c.owner, c.title, c.text, c.st, c.et)
		if c.err != nil {
			require.NotNil(t, err)
			require.Equal(t, c.err, err)
		}
	}
}
