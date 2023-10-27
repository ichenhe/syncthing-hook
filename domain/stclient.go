package domain

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/syncthing/syncthing/lib/events"
)

type SyncthingClient interface {
	// endpoint_system -------------------------------------------------------------------------

	GetSystemStatus() (*SystemStatus, error)

	// endpoint_events -------------------------------------------------------------------------

	// GetEvents receives events. since sets the ID of the last event you’ve already seen. The default value is 0, which
	// returns all events. The timeout duration can be customized with the parameter timeout(seconds).
	// To receive only a limited number of events, add the limit parameter with a suitable value for n and only the last n
	// events will be returned.
	//
	// Note: if more than one eventTypes filter is given, only subsequent events can be fetched. This function will wait
	// until a new event or timeout. However, if there's only one event type or empty, cached events can be fetched.
	// This is the intended behavior of Syncthing API, details: https://github.com/syncthing/syncthing/issues/8902
	GetEvents(eventTypes []events.EventType, since int, timeout int, limit int) ([]events.Event, error)

	GetDiskEvents(since int, timeout int, limit int) ([]events.Event, error)

	SubscribeEvent(eventTypes []events.EventType, since int) <-chan events.Event

	UnsubscribeEvent(eventCh <-chan events.Event)
}

type StApiError struct {
	HttpStatusCode int
	Err            error
}

var _ error = (*StApiError)(nil)

func (e StApiError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("failed to request syncthing api: %s", e.Err)
	} else {
		return fmt.Sprintf("syncthing api error [status=%d]", e.HttpStatusCode)
	}
}

func (e StApiError) Unwrap() error {
	return e.Err
}

func NewStApiReqError(err error) StApiError {
	return StApiError{
		Err: err,
	}
}

func NewStApiHttpError(resp *resty.Response) StApiError {
	return StApiError{
		HttpStatusCode: resp.StatusCode(),
	}
}
