package interfaces

import (
	"context"
	"github.com/rendau/my-otus/task8/scheduler/internal/domain/entities"
)

// Storage - is interface of storage
type Storage interface {
	EventList(ctx context.Context, filter *entities.EventListFilter) ([]*entities.Event, error)
}
