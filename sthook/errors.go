package sthook

type eventTypeError struct {
	msg string
}

func (e *eventTypeError) Error() string {
	return e.msg
}

func newEventTypeError(message string) *eventTypeError {
	return &eventTypeError{
		msg: message,
	}
}
