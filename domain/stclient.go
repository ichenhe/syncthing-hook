package domain

import (
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
