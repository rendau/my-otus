package usecases

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/errors"
	"time"
)

// Event - is type for event usecases
type Event struct {
	rUcs *Usecases
}

// CreateEvent - creates new event usecases
func CreateEvent(rUcs *Usecases) *Event {
	return &Event{
		rUcs: rUcs,
	}
}

func (ucs *Event) validate(ctx context.Context, event *entities.Event) error {
	if event.Owner == "" {
		return errors.ErrOwnerRequired
	}
	if event.Title == "" {
		return errors.ErrTitleRequired
	}
	if event.StartTime.IsZero() {
		return errors.ErrStartDateRequired
	}
	if event.StartTime.Before(time.Now()) {
		return errors.ErrIncorrectStartDate
	}
	if event.EndTime.IsZero() {
		return errors.ErrEndDateRequired
	}
	if event.EndTime.Before(event.StartTime) {
		return errors.ErrEndDateLTStartDate
	}
	overlapEventCount, err := ucs.rUcs.stg.EventListCount(ctx, &entities.EventListFilter{
		IDNE:        &event.ID,
		StartTimeLt: &event.StartTime,
		EndTimeGt:   &event.StartTime,
	})
	if err != nil {
		return err
	}
	if overlapEventCount > 0 {
		return errors.ErrOverlaping
	}
	return nil
}

// List - returns list of event
func (ucs *Event) List(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error) {
	events, err := ucs.rUcs.stg.EventList(ctx, filter)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to get list of events", "err", err.Error())
	}

	return events, err
}

// Create - creates event
func (ucs *Event) Create(ctx context.Context,
	owner, title, text string, startTime, endTime time.Time) (*entities.Event, error) {
	event := &entities.Event{
		ID:        0,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := ucs.validate(ctx, event)
	if err != nil {
		return nil, err
	}
	err = ucs.rUcs.stg.EventCreate(ctx, event)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to create event", "err", err.Error())
		return nil, err
	}
	return event, nil
}

// Get - retrieves event
func (ucs *Event) Get(ctx context.Context, id int64) (*entities.Event, error) {
	events, err := ucs.rUcs.stg.EventGet(ctx, id)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to get event", "err", err.Error())
	}
	return events, err
}

// Update - updates event
func (ucs *Event) Update(ctx context.Context, id int64,
	owner, title, text string, startTime time.Time, endTime time.Time) error {
	event := &entities.Event{
		ID:        id,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := ucs.validate(ctx, event)
	if err != nil {
		return err
	}
	err = ucs.rUcs.stg.EventUpdate(ctx, id, event)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to update event", "err", err.Error())
		return err
	}
	return nil
}

// Delete - deletes event
func (ucs *Event) Delete(ctx context.Context, id int64) error {
	err := ucs.rUcs.stg.EventDelete(ctx, id)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to delete event", "err", err.Error())
	}
	return err
}
