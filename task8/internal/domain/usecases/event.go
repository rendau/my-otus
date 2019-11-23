package usecases

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/errors"
	"github.com/satori/go.uuid"
	"time"
)

// Event - is type for event usecases
type Event struct {
	rUks *Usecases
}

// CreateEvent - creates new event usecases
func CreateEvent(rUks *Usecases) *Event {
	return &Event{
		rUks: rUks,
	}
}

func (uks *Event) validate(ctx context.Context, event *entities.Event) error {
	if event.Owner == "" {
		return errors.ErrOwnerRequired
	}
	if event.Title == "" {
		return errors.ErrTitleRequired
	}
	if event.StartTime.Before(time.Now()) {
		return errors.ErrIncorrectStartDate
	}
	if event.EndTime.Before(event.StartTime) {
		return errors.ErrEndDateLTStartDate
	}
	overlapEventCount, err := uks.rUks.stg.EventListCount(ctx, &entities.EventListFilter{
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
func (uks *Event) List(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error) {
	return uks.rUks.stg.EventList(ctx, filter)
}

// Create - creates event
func (uks *Event) Create(ctx context.Context,
	owner, title, text string, startTime time.Time, endTime time.Time) (*entities.Event, error) {
	uuidID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	event := &entities.Event{
		ID:        uuidID.String(),
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err = uks.validate(ctx, event)
	if err != nil {
		return nil, err
	}
	err = uks.rUks.stg.EventCreate(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Get - retrieves event
func (uks *Event) Get(ctx context.Context, id string) (*entities.Event, error) {
	return uks.rUks.stg.EventGet(ctx, id)
}

// Update - updates event
func (uks *Event) Update(ctx context.Context, id string,
	owner, title, text string, startTime time.Time, endTime time.Time) error {
	event := &entities.Event{
		ID:        id,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := uks.validate(ctx, event)
	if err != nil {
		return err
	}
	err = uks.rUks.stg.EventUpdate(ctx, id, event)
	if err != nil {
		return err
	}
	return nil
}

// Delete - deletes event
func (uks *Event) Delete(ctx context.Context, id string) error {
	return uks.rUks.stg.EventDelete(ctx, id)
}
