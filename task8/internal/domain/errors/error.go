package errors

type EventError string

func (e EventError) Error() string {
	return string(e)
}

var (
	ErrOverlaping       = EventError("another event exists for this date")
	ErrIncorrectEndDate = EventError("end_date is incorrect")
)
