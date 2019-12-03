package errors

// EventError - event error type
type EventError string

func (e EventError) Error() string {
	return string(e)
}

var (
	// ErrOverlaping - is error for overlaping
	ErrOverlaping = EventError("another event exists for this date")
	// ErrOwnerRequired - is error for "owner is required"
	ErrOwnerRequired = EventError("owner is required")
	// ErrTitleRequired - is error for "title is required"
	ErrTitleRequired = EventError("title is required")
	// ErrStartDateRequired - is error for "start date is required"
	ErrStartDateRequired = EventError("start date is required")
	// ErrEndDateRequired - is error for "end date is required"
	ErrEndDateRequired = EventError("end date is required")
	// ErrIncorrectStartDate - is error for "start_date is incorrect"
	ErrIncorrectStartDate = EventError("start_date is incorrect")
	// ErrEndDateLTStartDate - is error for "end_date is less than start_date"
	ErrEndDateLTStartDate = EventError("end_date is less than start_date")
)
