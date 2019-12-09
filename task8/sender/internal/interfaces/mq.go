package interfaces

import (
	"context"
	"github.com/rendau/my-otus/task8/sender/internal/domain/entities"
)

// Mq - is interface of mq
type Mq interface {
	GetEvent(ctx context.Context) (*entities.Event, error)
}
