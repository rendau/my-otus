package usecases

import (
	"github.com/rendau/my-otus/task8/internal/interfaces"
)

// Usecases - is root level usecases
type Usecases struct {
	log interfaces.Logger
	stg interfaces.Storage

	// modules
	Event *Event
}

// CreateUsecases - creates root level usecases instance
func CreateUsecases(log interfaces.Logger, stg interfaces.Storage) *Usecases {
	rUcs := &Usecases{
		log: log,
		stg: stg,
	}

	rUcs.Event = CreateEvent(rUcs)

	return rUcs
}
