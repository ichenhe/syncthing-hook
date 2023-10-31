package domain

import "github.com/syncthing/syncthing/lib/events"

type EventCh = <-chan *Event

// ConvertNativeEventChannel converts native event channel into custom event channel. The returned channel will be
// closed automatically once the given native event channel is closed.
//
// Events that failed to convert to local event will be ignored.
func ConvertNativeEventChannel(nativeEventCh <-chan events.Event) EventCh {
	ch := make(chan *Event)
	go func() {
		for nativeEv := range nativeEventCh {
			if ev, err := NewEventFromStEvent(nativeEv); err != nil {
				continue
			} else {
				ch <- ev
			}
		}
		close(ch)
	}()
	return ch
}

// EventSource is a unified event subscription portal, including native syncthing event and extra event detected by
// syncthing hook.
type EventSource interface {
	Subscribe(eventType EventType, params *HookParameters, hookDef *HookDefinition) (EventCh, error)
	// Unsubscribe must eventually close the given channel.
	Unsubscribe(eventCh EventCh)
}

type IllegalEventParamError struct {
	Message string
}

func NewIllegalEventParamError(message string) *IllegalEventParamError {
	return &IllegalEventParamError{
		Message: message,
	}
}

func (e *IllegalEventParamError) Error() string {
	return e.Message
}
