package domain

import (
	"errors"
	"fmt"
	"github.com/syncthing/syncthing/lib/events"
	"strings"
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

func NewEventFromStEvent(event events.Event) *Event {
	return &Event{
		Time: event.Time,
		Type: convertFromStEventType(event.Type),
		Data: event.Data,
	}
}

func convertFromStEventType(eventType events.EventType) EventType {
	switch eventType {
	case events.FolderCompletion:
		return FolderCompletion
	case events.LocalChangeDetected:
		return LocalChangeDetected
	}
	panic(errors.New(fmt.Sprintf("unknown st exevent type '%d'", eventType)))
}

func NewEvent(eventType EventType, eventTime time.Time, data any) *Event {
	return &Event{
		Time: eventTime,
		Type: eventType,
		Data: data,
	}
}

type EventType string

var ErrNotValidNativeEventType = errors.New("not a valid native event type")

const (
	UnknownEventType    EventType = ""
	FolderCompletion    EventType = "st:FolderCompletion"
	LocalChangeDetected EventType = "st:LocalChangeDetected"

	LocalFolderContentChangeDetected EventType = "ex:LocalFolderContentChangeDetected"
)

func UnmarshalEventType(evType string) EventType {
	switch evType {
	case string(FolderCompletion):
		return FolderCompletion
	case string(LocalFolderContentChangeDetected):
		return LocalFolderContentChangeDetected
	default:
		return UnknownEventType
	}
}

func (t EventType) IsNativeEvent() bool {
	return strings.HasPrefix(string(t), "st:")
}

// ConvertToNative converts current event type to syncthing native type. Returns ErrNotValidNativeEventType if failed.
func (t EventType) ConvertToNative() (events.EventType, error) {
	if !t.IsNativeEvent() {
		return 0, ErrNotValidNativeEventType
	}
	native := events.UnmarshalEventType(string(t)[3:])
	if native == 0 {
		return 0, ErrNotValidNativeEventType
	} else {
		return native, nil
	}
}
