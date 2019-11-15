package memdb

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
)

func (mdb *MemDb) EventList(ctx context.Context) []*entities.Event {
	// TODO
	return nil
}

func (mdb *MemDb) EventCreate(ctx context.Context, event *entities.Event) error {
	// TODO
	return nil
}

func (mdb *MemDb) EventGet(ctx context.Context, id string) (*entities.Event, error) {
	// TODO
	return nil, nil
}

func (mdb *MemDb) EventUpdate(ctx context.Context, id string, event *entities.Event) error {
	// TODO
	return nil
}

func (mdb *MemDb) EventDelete(ctx context.Context, id string) (*entities.Event, error) {
	// TODO
	return nil, nil
}
