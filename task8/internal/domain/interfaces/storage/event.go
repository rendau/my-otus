package storage

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
)

// Event - is interface of event for db-adapters
type Event interface {
	EventList(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error)
	EventListCount(ctx context.Context, filter *entities.EventListFilter) (int64, error)
	EventCreate(ctx context.Context, event *entities.Event) error
	EventGet(ctx context.Context, id string) (*entities.Event, error)
	EventUpdate(ctx context.Context, id string, event *entities.Event) error
	EventDelete(ctx context.Context, id string) error
}
