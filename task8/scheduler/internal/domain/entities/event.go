package entities

import (
	"time"
)

// Event - is type for event
type Event struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// EventListFilter - is type for event-filter
type EventListFilter struct {
	StartTimeLt *time.Time // if start time less than value
	StartTimeGt *time.Time // if start time greater than value
	EndTimeLt   *time.Time // if end time less than value
	EndTimeGt   *time.Time // if end time greater than value
}
