package internal

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

// EventCE - is type for event create/edit
type EventCE struct {
	Owner     *string    `json:"owner"`
	Title     *string    `json:"title"`
	Text      *string    `json:"text"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}
