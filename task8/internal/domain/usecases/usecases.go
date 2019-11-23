package usecases

import (
	"github.com/rendau/my-otus/task8/internal/config"
	"github.com/rendau/my-otus/task8/internal/interfaces"
)

// Usecases - is root level usecases
type Usecases struct {
	cfg *config.Config
	stg interfaces.Storage

	// modules
	Event *Event
}

// CreateUsecases - creates root level usecases instance
func CreateUsecases(cfg *config.Config, stg interfaces.Storage) *Usecases {
	rUks := &Usecases{
		cfg: cfg,
		stg: stg,
	}

	rUks.Event = CreateEvent(rUks)

	return rUks
}
