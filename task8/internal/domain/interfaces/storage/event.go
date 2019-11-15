package storage

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
)

type Event interface {
	List(ctx context.Context) []*entities.Event
	Create(ctx context.Context, event *entities.Event) error
	Get(ctx context.Context, id string) (*entities.Event, error)
	Update(ctx context.Context, id string, event *entities.Event) error
	Delete(ctx context.Context, id string) (*entities.Event, error)
}
