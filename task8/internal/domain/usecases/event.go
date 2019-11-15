package usecases

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/interfaces/storage"
	"github.com/satori/go.uuid"
	"time"
)

type EventUsecases struct {
	stg storage.Event
}

func (es *EventUsecases) CreateEvent(ctx context.Context, owner, title, text string, startTime, endTime *time.Time) (*entities.Event, error) {
	// TODO: persistence, validation
	event := &entities.Event{
		Id:        uuid.Must(uuid.NewV4()),
		Owner:     owner,
		Title:     title,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := es.stg.Create(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}
