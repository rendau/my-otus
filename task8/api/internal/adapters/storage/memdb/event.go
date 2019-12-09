package memdb

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/domain/entities"
	"sync"
)

type eventTableSt struct {
	t     []*entities.Event
	idSeq int64
	mu    sync.RWMutex
}

// EventList - list events
func (mdb *MemDb) EventList(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error) {
	mdb.eventTable.mu.RLock()
	defer mdb.eventTable.mu.RUnlock()
	res := make([]*entities.Event, 0, 20)
	for _, e := range mdb.eventTable.t {
		if filter == nil {
			res = append(res, e)
			continue
		}
		if filter.IDNE != nil && e.ID == *filter.IDNE {
			continue
		}
		if filter.StartTimeLt != nil && !e.StartTime.Before(*filter.StartTimeLt) {
			continue
		}
		if filter.StartTimeGt != nil && !e.StartTime.After(*filter.StartTimeGt) {
			continue
		}
		if filter.EndTimeLt != nil && !e.EndTime.Before(*filter.EndTimeLt) {
			continue
		}
		if filter.EndTimeGt != nil && !e.EndTime.After(*filter.EndTimeGt) {
			continue
		}
		res = append(res, e)
	}
	return res, nil
}

// EventListCount - count of filtered list
func (mdb *MemDb) EventListCount(ctx context.Context, filter *entities.EventListFilter) (int64, error) {
	events, err := mdb.EventList(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int64(len(events)), nil
}

// EventCreate - creates event
func (mdb *MemDb) EventCreate(ctx context.Context, event *entities.Event) error {
	mdb.eventTable.mu.Lock()
	defer mdb.eventTable.mu.Unlock()

	mdb.eventTable.idSeq++

	dbEvent := &entities.Event{
		ID:        mdb.eventTable.idSeq,
		Owner:     event.Owner,
		Title:     event.Title,
		Text:      event.Text,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
	}
	mdb.eventTable.t = append(mdb.eventTable.t, dbEvent)

	event.ID = dbEvent.ID

	return nil
}

// EventGet - retrieves one event
func (mdb *MemDb) EventGet(ctx context.Context, id int64) (*entities.Event, error) {
	mdb.eventTable.mu.RLock()
	defer mdb.eventTable.mu.RUnlock()
	for _, e := range mdb.eventTable.t {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, nil
}

// EventUpdate - updates event by id
func (mdb *MemDb) EventUpdate(ctx context.Context, id int64, event *entities.Event) error {
	mdb.eventTable.mu.Lock()
	defer mdb.eventTable.mu.Unlock()
	for _, e := range mdb.eventTable.t {
		if e.ID != id {
			continue
		}
		e.Owner = event.Owner
		e.Title = event.Title
		e.Text = event.Text
		e.StartTime = event.StartTime
		e.EndTime = event.EndTime
		break
	}
	return nil
}

// EventDelete - deletes event by id
func (mdb *MemDb) EventDelete(ctx context.Context, id int64) error {
	mdb.eventTable.mu.Lock()
	defer mdb.eventTable.mu.Unlock()
	newT := make([]*entities.Event, 0, len(mdb.eventTable.t))
	for _, e := range mdb.eventTable.t {
		if e.ID != id {
			newT = append(newT, e)
		}
	}
	mdb.eventTable.t = newT
	return nil
}
