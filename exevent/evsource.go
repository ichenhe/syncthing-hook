package exevent

import (
	"SyncthingHook/domain"
	"fmt"
	"github.com/syncthing/syncthing/lib/events"
	"go.uber.org/zap"
	"sync"
)

type eventUnsubscriber interface {
	Unsubscribe()
}

type eventUnsubscriberFunc func()

func (f eventUnsubscriberFunc) Unsubscribe() {
	f()
}

type EventSource struct {
	stClient domain.SyncthingClient
	logger   *zap.SugaredLogger

	locker        sync.RWMutex
	unsubscribers map[domain.EventCh]eventUnsubscriber
}

var _ domain.EventSource = (*EventSource)(nil)

func NewEventSource(stClient domain.SyncthingClient, logger *zap.SugaredLogger) domain.EventSource {
	return &EventSource{
		stClient:      stClient,
		logger:        logger,
		unsubscribers: make(map[domain.EventCh]eventUnsubscriber),
	}
}

// Subscribe subscribes event of given type and apply filters based on params.
//
// Returned error could be domain.IllegalEventParamError | domain.ErrNotValidNativeEventType
func (s *EventSource) Subscribe(eventType domain.EventType, params *domain.HookParameters, hookDef *domain.HookDefinition) (domain.EventCh, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	logger := s.logger
	if hookDef != nil {
		logger = hookDef.AddToLogger(logger)
	}

	var eventCh domain.EventCh
	var eventUnsub eventUnsubscriber
	if eventType.IsNativeEvent() {
		nativeType, err := eventType.ConvertToNative()
		if err != nil {
			return nil, err
		}
		nativeCh := s.stClient.SubscribeEvent([]events.EventType{nativeType}, 0)
		eventCh = domain.ConvertNativeEventChannel(nativeCh)
		eventUnsub = eventUnsubscriberFunc(func() {
			s.stClient.UnsubscribeEvent(nativeCh)
		})
	} else {
		switch eventType {
		case domain.LocalFolderContentChangeDetected:
			if ch, unsub, err := detectLocalFolderContentChanged(s.stClient, params, logger); err != nil {
				return nil, err
			} else {
				eventCh = ch
				eventUnsub = unsub
			}
		default:
			return nil, domain.NewIllegalEventParamError(fmt.Sprintf("unknown eventType: %s", eventType))
		}
	}

	// add filters
	filters := s.createFilters(params, logger)
	eventCh = applyHandlers(eventCh, filters...)

	s.unsubscribers[eventCh] = eventUnsub
	return eventCh, nil
}

// createFilters creates filter-type event handlers. Returns nil or empty if given params is nil or no filter needed.
func (s *EventSource) createFilters(params *domain.HookParameters, logger *zap.SugaredLogger) []domain.EventHandler {
	if params == nil {
		return nil
	}
	handlers := make([]domain.EventHandler, 0)
	if params.ContainsKey("tolerance") {
		handlers = append(handlers, NewToleranceFilter(params.GetInt64("tolerance", 0), logger))
	}
	if params.ContainsKey("cooldown") {
		handlers = append(handlers, NewCoolDownFilter(params.GetInt64("cooldown", 0), logger))
	}
	return handlers
}

func (s *EventSource) Unsubscribe(eventCh domain.EventCh) {
	s.locker.Lock()
	defer s.locker.Unlock()

	f, ex := s.unsubscribers[eventCh]
	if !ex {
		s.logger.Info("given eventCh not exist, ignore unsubscribe request")
		return
	}
	delete(s.unsubscribers, eventCh)
	f.Unsubscribe()
}
