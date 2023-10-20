package extevent

import (
	"github.com/syncthing/syncthing/lib/events"
	"time"
)

// eventTypes creates a slice of events.EventType for a single type to simplify the call.
func eventTypes(eventType events.EventType) []events.EventType {
	return []events.EventType{eventType}
}

func stopTimerIfNeeded(timer *time.Timer) {
	if timer != nil && !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}
