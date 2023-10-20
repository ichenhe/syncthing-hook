package extevent

type IllegalArgumentError struct {
	Message string
}

func (i IllegalArgumentError) Error() string {
	return i.Message
}

func newIllegalArgumentError(message string) *IllegalArgumentError {
	return &IllegalArgumentError{Message: message}
}
