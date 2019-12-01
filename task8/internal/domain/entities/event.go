package entities

import (
	"time"
)

// Event - is type for event
type Event struct {
	ID        int64
	Owner     string
	Title     string
	Text      string
	StartTime time.Time
	EndTime   time.Time
}

// EventListFilter - is type for event-filter
type EventListFilter struct {
	StartTimeLt *time.Time // if start time less than value
	StartTimeGt *time.Time // if start time greater than value
	EndTimeLt   *time.Time // if end time less than value
	EndTimeGt   *time.Time // if end time greater than value
}
