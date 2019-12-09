package interfaces

import (
	"github.com/rendau/my-otus/task8/scheduler/internal/domain/entities"
)

// Mq - is interface of mq
type Mq interface {
	PublishEventNotification(event *entities.Event) error
}
