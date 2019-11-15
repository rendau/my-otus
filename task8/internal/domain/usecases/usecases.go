package usecases

import (
	"github.com/rendau/my-otus/task8/internal/domain/interfaces/storage"
)

// Usecases - type for root usecases
type Usecases struct {
	stg storage.Event

	// modules
	Event *Event
}

// CreateUsecases - creates root usecases instance
func CreateUsecases(stg storage.Event) *Usecases {
	return &Usecases{
		stg: stg,

		// modules
		Event: NewEvent(stg),
	}
}
