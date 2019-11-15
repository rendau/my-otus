package usecases

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/errors"
	"github.com/rendau/my-otus/task8/internal/domain/interfaces/storage"
	"github.com/satori/go.uuid"
	"time"
)

// Event - is type for event usecases
type Event struct {
	stg storage.Event
}

// NewEvent - creates new event usecases
func NewEvent(stg storage.Event) *Event {
	return &Event{
		stg: stg,
	}
}

func (uk *Event) validate(ctx context.Context, event *entities.Event) error {
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
	overlapEventCount, err := uk.stg.EventListCount(ctx, &entities.EventListFilter{
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
func (uk *Event) List(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error) {
	return uk.stg.EventList(ctx, filter)
}

// Create - creates event
func (uk *Event) Create(ctx context.Context,
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
	err = uk.validate(ctx, event)
	if err != nil {
		return nil, err
	}
	err = uk.stg.EventCreate(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Get - retrieves event
func (uk *Event) Get(ctx context.Context, id string) (*entities.Event, error) {
	return uk.stg.EventGet(ctx, id)
}

// Update - updates event
func (uk *Event) Update(ctx context.Context, id string,
	owner, title, text string, startTime time.Time, endTime time.Time) error {
	event := &entities.Event{
		ID:        id,
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := uk.validate(ctx, event)
	if err != nil {
		return err
	}
	err = uk.stg.EventUpdate(ctx, id, event)
	if err != nil {
		return err
	}
	return nil
}

// Delete - deletes event
func (uk *Event) Delete(ctx context.Context, id string) error {
	return uk.stg.EventDelete(ctx, id)
}
