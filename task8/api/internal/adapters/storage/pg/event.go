package pg

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/domain/entities"
	"time"
)

// EventList - list events
func (pdb *PostgresDb) EventList(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error) {
	qFrom := ` from event e`

	cond, dbArgs := pdb.eventListCond("e", filter)

	qWhere := ` where 1=1 ` + cond

	qSelect := ` select e.id, e.owner, e.title, e.text, e.start_time, e.end_time`

	qOrderBy := ` order by e.start_time, e.end_time`

	stmt, err := pdb.db.PrepareNamed(qSelect + qFrom + qWhere + qOrderBy)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, dbArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entities.Event

	for rows.Next() {
		event := &entities.Event{}
		err = rows.Scan(&event.ID, &event.Owner, &event.Title, &event.Text, &event.StartTime, &event.EndTime)
		if err != nil {
			return nil, err
		}
		event.StartTime = event.StartTime.Local()
		event.EndTime = event.EndTime.Local()
		events = append(events, event)
	}

	return events, nil
}

// EventListCount - count of filtered list
func (pdb *PostgresDb) EventListCount(ctx context.Context, filter *entities.EventListFilter) (int64, error) {
	qFrom := ` from event e`

	cond, dbArgs := pdb.eventListCond("e", filter)

	qWhere := ` where 1=1 ` + cond

	qSelect := ` select count(*)`

	stmt, err := pdb.db.PrepareNamed(qSelect + qFrom + qWhere)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var res int64

	err = stmt.GetContext(ctx, &res, dbArgs)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (pdb *PostgresDb) eventListCond(tAlias string, filter *entities.EventListFilter) (string, map[string]interface{}) {
	cond := ``
	dbArgs := map[string]interface{}{}

	if filter != nil {
		if filter.IDNE != nil {
			cond += ` and ` + tAlias + `.id != :id_ne`
			dbArgs["id_ne"] = *filter.IDNE
		}

		if filter.StartTimeLt != nil {
			cond += ` and ` + tAlias + `.start_time < :start_time_lt`
			dbArgs["start_time_lt"] = (*filter.StartTimeLt).Format(time.RFC3339)
		}

		if filter.StartTimeGt != nil {
			cond += ` and ` + tAlias + `.start_time > :start_time_gt`
			dbArgs["start_time_gt"] = (*filter.StartTimeGt).Format(time.RFC3339)
		}

		if filter.EndTimeLt != nil {
			cond += ` and ` + tAlias + `.end_time < :end_time_lt`
			dbArgs["end_time_lt"] = (*filter.EndTimeLt).Format(time.RFC3339)
		}

		if filter.EndTimeGt != nil {
			cond += ` and ` + tAlias + `.end_time > :end_time_gt`
			dbArgs["end_time_gt"] = (*filter.EndTimeGt).Format(time.RFC3339)
		}
	}

	return cond, dbArgs
}

// EventCreate - creates event
func (pdb *PostgresDb) EventCreate(ctx context.Context, event *entities.Event) error {
	stmt, err := pdb.db.PrepareNamed(`
		insert into event(owner, title, text, start_time, end_time)
		values (:owner, :title, :text, :start_time, :end_time)
		returning id
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &event.ID, map[string]interface{}{
		"owner":      event.Owner,
		"title":      event.Title,
		"text":       event.Text,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
	})
	if err != nil {
		return err
	}

	return nil
}

// EventGet - retrieves one event
func (pdb *PostgresDb) EventGet(ctx context.Context, id int64) (*entities.Event, error) {
	event := &entities.Event{}

	err := pdb.db.QueryRowxContext(ctx, `
		select id, owner, title, text, start_time, end_time
		from event
		where id = $1
	`, id).Scan(&event.ID, &event.Owner, &event.Title, &event.Text, &event.StartTime, &event.EndTime)

	event.StartTime = event.StartTime.Local()
	event.EndTime = event.EndTime.Local()

	return event, err
}

// EventUpdate - updates event by id
func (pdb *PostgresDb) EventUpdate(ctx context.Context, id int64, event *entities.Event) error {
	stmt, err := pdb.db.PrepareNamed(`
		update event
		    set owner=:owner, title=:title, text=:text, start_time=:start_time, end_time=:end_time
		where id = :id
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, map[string]interface{}{
		"id":         id,
		"owner":      event.Owner,
		"title":      event.Title,
		"text":       event.Text,
		"start_time": event.StartTime,
		"end_time":   event.EndTime,
	})
	if err != nil {
		return err
	}

	return nil
}

// EventDelete - deletes event by id
func (pdb *PostgresDb) EventDelete(ctx context.Context, id int64) error {
	_, err := pdb.db.ExecContext(ctx, `delete from event where id = $1`, id)
	return err
}
