package pg

import (
	"context"
	"github.com/rendau/my-otus/task8/scheduler/internal/domain/entities"
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

func (pdb *PostgresDb) eventListCond(tAlias string, filter *entities.EventListFilter) (string, map[string]interface{}) {
	cond := ``
	dbArgs := map[string]interface{}{}

	if filter != nil {
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
