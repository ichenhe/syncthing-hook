package domain

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

func NewEventFromStEvent(event events.Event) (*Event, error) {
	ev := convertFromStEventType(event.Type)
	if ev == UnknownEventType {
		return nil, errors.New("unknown syncthing event type")
	}
	return &Event{
		Time: event.Time,
		Type: ev,
		Data: event.Data,
	}, nil
}

// convertFromStEventType converts a Syncthing event type to local one. Returns UnknownEventType if error.
func convertFromStEventType(eventType events.EventType) EventType {
	if locEv, ex := mapToLocalType[eventType]; ex {
		return locEv
	}
	return UnknownEventType
}

type EventType string

var ErrNotValidNativeEventType = errors.New("not a valid native event type")

// Supported event types. Each type except UnknownEventType MUST be registered in mapToStType.
const (
	UnknownEventType EventType = ""

	ConfigSaved             EventType = "st:ConfigSaved"
	DeviceConnected         EventType = "st:DeviceConnected"
	DeviceDisconnected      EventType = "st:DeviceDisconnected"
	DeviceDiscovered        EventType = "st:DeviceDiscovered"
	DevicePaused            EventType = "st:DevicePaused"
	DeviceResumed           EventType = "st:DeviceResumed"
	DownloadProgress        EventType = "st:DownloadProgress"
	Failure                 EventType = "st:Failure"
	FolderCompletion        EventType = "st:FolderCompletion"
	FolderErrors            EventType = "st:FolderErrors"
	FolderPaused            EventType = "st:FolderPaused"
	FolderResumed           EventType = "st:FolderResumed"
	FolderScanProgress      EventType = "st:FolderScanProgress"
	FolderSummary           EventType = "st:FolderSummary"
	FolderWatchStateChanged EventType = "st:FolderWatchStateChanged"
	ItemFinished            EventType = "st:ItemFinished"
	ItemStarted             EventType = "st:ItemStarted"
	ListenAddressesChanged  EventType = "st:ListenAddressesChanged"
	LocalChangeDetected     EventType = "st:LocalChangeDetected"
	LocalIndexUpdated       EventType = "st:LocalIndexUpdated"
	LoginAttempt            EventType = "st:LoginAttempt"
	PendingDevicesChanged   EventType = "st:PendingDevicesChanged"
	PendingFoldersChanged   EventType = "st:PendingFoldersChanged"
	RemoteChangeDetected    EventType = "st:RemoteChangeDetected"
	RemoteDownloadProgress  EventType = "st:RemoteDownloadProgress"
	RemoteIndexUpdated      EventType = "st:RemoteIndexUpdated"
	Starting                EventType = "st:Starting"
	StartupComplete         EventType = "st:StartupComplete"
	StateChanged            EventType = "st:StateChanged"

	LocalFolderContentChangeDetected EventType = "ex:LocalFolderContentChangeDetected"
)

// mapToStType records ALL events and corresponding native events, 0 if no corresponding event type.
var mapToStType = map[EventType]events.EventType{
	ConfigSaved:             events.ConfigSaved,
	DeviceConnected:         events.DeviceConnected,
	DeviceDisconnected:      events.DeviceDisconnected,
	DeviceDiscovered:        events.DeviceDiscovered,
	DevicePaused:            events.DevicePaused,
	DeviceResumed:           events.DeviceResumed,
	DownloadProgress:        events.DownloadProgress,
	Failure:                 events.Failure,
	FolderCompletion:        events.FolderCompletion,
	FolderErrors:            events.FolderErrors,
	FolderPaused:            events.FolderPaused,
	FolderResumed:           events.FolderResumed,
	FolderScanProgress:      events.FolderScanProgress,
	FolderSummary:           events.FolderSummary,
	FolderWatchStateChanged: events.FolderWatchStateChanged,
	ItemFinished:            events.ItemFinished,
	ItemStarted:             events.ItemStarted,
	ListenAddressesChanged:  events.ListenAddressesChanged,
	LocalChangeDetected:     events.LocalChangeDetected,
	LocalIndexUpdated:       events.LocalIndexUpdated,
	LoginAttempt:            events.LoginAttempt,
	PendingDevicesChanged:   events.PendingDevicesChanged,
	PendingFoldersChanged:   events.PendingFoldersChanged,
	RemoteChangeDetected:    events.RemoteChangeDetected,
	RemoteDownloadProgress:  events.RemoteDownloadProgress,
	RemoteIndexUpdated:      events.RemoteIndexUpdated,
	Starting:                events.Starting,
	StartupComplete:         events.StartupComplete,
	StateChanged:            events.StateChanged,

	LocalFolderContentChangeDetected: events.EventType(0),
}

var mapToLocalType map[events.EventType]EventType

func init() {
	// build reverse map
	mapToLocalType = make(map[events.EventType]EventType, len(mapToStType))
	for cu, st := range mapToStType {
		if st != 0 {
			mapToLocalType[st] = cu
		}
	}
}

func UnmarshalEventType(evType string) EventType {
	t := EventType(evType)
	if _, ex := mapToStType[t]; !ex {
		return UnknownEventType
	}
	return t
}

func (t EventType) IsNativeEvent() bool {
	e, ex := mapToStType[t]
	return ex && e != 0
}

// ConvertToNative converts current event type to syncthing native type. Returns ErrNotValidNativeEventType if failed.
func (t EventType) ConvertToNative() (events.EventType, error) {
	n, ex := mapToStType[t]
	if !ex || n == 0 {
		return 0, ErrNotValidNativeEventType
	}
	return n, nil
}
