package extevent

import (
	"errors"
	"fmt"
	"github.com/syncthing/syncthing/lib/events"
	"time"
)

type Event struct {
	Time time.Time   `json:"time"`
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

func (e *Event) String() string {
	return fmt.Sprintf("Event <time: %s, type: %s>", e.Time.String(), e.Type)
}

func NewEventFromStEvent(event events.Event) Event {
	return Event{
		Time: event.Time,
		Type: convertFromStEventType(event.Type),
		Data: event.Data,
	}
}

func NewEvent(eventType EventType, eventTime time.Time, data any) *Event {
	return &Event{
		Time: eventTime,
		Type: eventType,
		Data: data,
	}
}

type EventType = string

const (
	FolderCompletion EventType = "st:FolderCompletion"

	LocalFolderContentChangeDetected EventType = "ex:LocalFolderContentChangeDetected"
)

func convertFromStEventType(eventType events.EventType) EventType {
	switch eventType {
	case events.FolderCompletion:
		return FolderCompletion
	}
	panic(errors.New(fmt.Sprintf("unknown st event type '%d'", eventType)))
}
