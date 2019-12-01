package interfaces

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
)

// Storage - is interface of storage
type Storage interface {
	EventList(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error)
	EventListCount(ctx context.Context, filter *entities.EventListFilter) (int64, error)
	EventCreate(ctx context.Context, event *entities.Event) error
	EventGet(ctx context.Context, id int64) (*entities.Event, error)
	EventUpdate(ctx context.Context, id int64, event *entities.Event) error
	EventDelete(ctx context.Context, id int64) error
}
