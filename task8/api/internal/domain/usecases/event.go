package usecases

import (
	"context"
	"github.com/rendau/my-otus/task8/api/internal/domain/entities"
	"github.com/rendau/my-otus/task8/api/internal/domain/errors"
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
		ucs.rUcs.log.Errorw("Fail to get list of events", "error", err.Error())
	}

	return events, err
}

// Create - creates event
func (ucs *Event) Create(ctx context.Context, event *entities.Event) error {
	event.ID = 0
	err := ucs.validate(ctx, event)
	if err != nil {
		return err
	}
	err = ucs.rUcs.stg.EventCreate(ctx, event)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to create event", "error", err.Error())
		return err
	}
	return nil
}

// Get - retrieves event
func (ucs *Event) Get(ctx context.Context, id int64) (*entities.Event, error) {
	events, err := ucs.rUcs.stg.EventGet(ctx, id)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to get event", "error", err.Error())
	}
	return events, err
}

// Update - updates event
func (ucs *Event) Update(ctx context.Context, event *entities.Event) error {
	err := ucs.validate(ctx, event)
	if err != nil {
		return err
	}
	err = ucs.rUcs.stg.EventUpdate(ctx, event.ID, event)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to update event", "error", err.Error())
		return err
	}
	return nil
}

// Delete - deletes event
func (ucs *Event) Delete(ctx context.Context, id int64) error {
	err := ucs.rUcs.stg.EventDelete(ctx, id)
	if err != nil {
		ucs.rUcs.log.Errorw("Fail to delete event", "error", err.Error())
	}
	return err
}
