package stclient

import (
	"SyncthingHook/safechan"
	"github.com/syncthing/syncthing/lib/events"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// GetEvents receives events. since sets the ID of the last event you’ve already seen. The default value is 0, which
// returns all events. The timeout duration can be customized with the parameter timeout(seconds).
// To receive only a limited number of events, add the limit parameter with a suitable value for n and only the last n
// events will be returned.
//
// Note: if more than one eventTypes filter is given, only subsequent events can be fetched. This function will wait
// until a new event or timeout. However, if there's only one event type or empty, cached events can be fetched.
// This is the intended behavior of Syncthing API, details: https://github.com/syncthing/syncthing/issues/8902
func (s *Syncthing) GetEvents(eventTypes []events.EventType, since int, timeout int, limit int) ([]events.Event, error) {
	types := make([]string, len(eventTypes))
	for i, e := range eventTypes {
		types[i] = e.String()
	}
	return s.getEventsWithStringTypes(strings.Join(types, ","), since, timeout, limit)
}

func (s *Syncthing) getEventsWithStringTypes(eventTypes string, since int, timeout int, limit int) ([]events.Event, error) {
	var result []events.Event
	params := map[string]string{
		"events":  eventTypes,
		"since":   strconv.Itoa(since),
		"timeout": strconv.Itoa(timeout),
		"limit":   strconv.Itoa(limit),
	}
	if resp, err := s.newRequest(result).SetQueryParams(params).Get("/rest/events"); err != nil {
		return nil, newApiError(err)
	} else if resp.IsError() {
		return nil, newHttpApiError(resp)
	} else {
		return *resp.Result().(*[]events.Event), nil
	}
}

func (s *Syncthing) GetDiskEvents(since int, timeout int, limit int) ([]events.Event, error) {
	var result []events.Event
	params := map[string]string{
		"since":   strconv.Itoa(since),
		"timeout": strconv.Itoa(timeout),
		"limit":   strconv.Itoa(limit),
	}
	if resp, err := s.newRequest(result).SetQueryParams(params).Get("/rest/events/disk"); err != nil {
		return nil, newApiError(err)
	} else if resp.IsError() {
		return nil, newHttpApiError(resp)
	} else {
		return *resp.Result().(*[]events.Event), nil
	}
}

type subscriber = *safechan.SafeChannel[events.Event]

type eventSubscription struct {
	downstreamLocker sync.RWMutex
	subscribers      map[events.EventType][]subscriber  // To dispatch different types of events to subscribers
	mapping          map[<-chan events.Event]subscriber // To locate subscriber to close it via channel returned to the caller

	upstreamLocker sync.RWMutex
	upstream       subscriber
}

func newEventSubscription() *eventSubscription {
	return &eventSubscription{
		subscribers: make(map[events.EventType][]subscriber),
		mapping:     make(map[<-chan events.Event]subscriber),
	}
}

func (s *Syncthing) UnsubscribeEvent(eventCh <-chan events.Event) {
	if eventCh == nil {
		return
	}
	s.eventSubscription.downstreamLocker.Lock()
	defer s.eventSubscription.downstreamLocker.Unlock()

	realCh, ex := s.eventSubscription.mapping[eventCh]
	if !ex {
		s.logger.Info("given subscriber not exist, ignore unsubscribe request")
		return // already unsubscribed
	}
	// remove the subscriber from list
	for eventType, subscribers := range s.eventSubscription.subscribers {
		for i, ch := range subscribers {
			if ch == realCh {
				newSubscribers := append(subscribers[:i], subscribers[i+1:]...)
				if len(newSubscribers) > 0 {
					s.eventSubscription.subscribers[eventType] = newSubscribers
				} else {
					s.logger.Debugw("subscribers list is empty, delete this type item", "eventType", eventType)
					delete(s.eventSubscription.subscribers, eventType)
				}
				break // don't return here because channel may exist in different event type
			}
		}
	}
	// close subscriber channel
	realCh.SafeClose()
	// delete mapping between channel returned and real channel
	delete(s.eventSubscription.mapping, eventCh)
	s.logger.Debugw("downstream unsubscribed", "remaining", len(s.eventSubscription.mapping))
	// unsubscribe universal subscriber if no long needed
	if len(s.eventSubscription.mapping) == 0 {
		s.logger.Debug("unsubscribe upstream because there's no downstream subscriber any more")
		s.eventSubscription.upstreamLocker.Lock()
		defer s.eventSubscription.upstreamLocker.Unlock()
		s.eventSubscription.upstream.SafeClose()
		s.eventSubscription.upstream = nil
	}
}

// SubscribeEvent subscribes any events compliant with eventTypes after since. Returns the event callback. Call
// unsubscribe if you want to stop subscription.
func (s *Syncthing) SubscribeEvent(eventTypes []events.EventType, since int) <-chan events.Event {
	s.eventSubscription.downstreamLocker.Lock()
	defer s.eventSubscription.downstreamLocker.Unlock()

	newSubscriber := safechan.NewSafeChannel[events.Event]()
	outSubscriber := newSubscriber.OutCh()
	for _, eventType := range eventTypes {
		// add downstream subscriber
		if subscribers, ex := s.eventSubscription.subscribers[eventType]; ex {
			s.logger.Debugw("append subscriber to list", "eventType", eventType)
			s.eventSubscription.subscribers[eventType] = append(subscribers, newSubscriber)
		} else {
			s.logger.Debugw("create subscriber list and attach the new item", "eventType", eventType)
			s.eventSubscription.subscribers[eventType] = []subscriber{newSubscriber}
		}
	}
	s.eventSubscription.mapping[outSubscriber] = newSubscriber

	// register upstream subscriber if needed
	if len(s.eventSubscription.mapping) == 1 {
		s.startSubscription(since)
	}
	return outSubscriber
}

// startSubscription starts polling for all event types. The event will be sent to
// s.eventSubscription.upstream. This function must not be called if subscription was already stated.
func (s *Syncthing) startSubscription(since int) {
	s.logger.Debug("subscribe upstream events")
	s.eventSubscription.upstreamLocker.Lock()
	s.eventSubscription.upstream = safechan.NewSafeChannel[events.Event]()
	s.eventSubscription.upstreamLocker.Unlock()

	// start listening and forwarding event to subscribers
	go func() {
		s.logger.Debug("start consuming upstream events")
		s.eventSubscription.upstreamLocker.RLock()
		upstream := s.eventSubscription.upstream
		s.eventSubscription.upstreamLocker.RUnlock()
		if upstream != nil { // user may unsubscribe before async subscribe process completed
			for event := range upstream.OutCh() {
				s.dispatchEvent(event)
			}
		}
	}()

	realSince := since
	const timeout int = 60
	s.logger.Debugw("start polling upstream events", "since", realSince)
	go func() {
		for { // polling
			resp := s.getAllEvents(realSince, timeout)
			if len(resp) == 0 {
				continue
			}

			s.eventSubscription.upstreamLocker.RLock()
			universalSubscriber := s.eventSubscription.upstream
			s.eventSubscription.upstreamLocker.RUnlock()
			if universalSubscriber == nil {
				return // user may unsubscribe before async subscribe process completed
			}
			for _, ev := range resp {
				ok := universalSubscriber.SafeSend(ev)
				if !ok {
					s.logger.Debug("universal subscriber closed, stop polling")
					return // subscribe channel closed
				}
				realSince = ev.SubscriptionID
			}
		}
	}()
}

func (s *Syncthing) getAllEvents(since int, timeout int) []events.Event {
	resp, _ := s.getEventsWithStringTypes(s.allEventTypes, since, timeout, 0)
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].SubscriptionID < resp[j].SubscriptionID
	})
	return resp
}

// dispatchEvent dispatches event from universal channel to subscribers. It requires downstreamLocker's read locker.
func (s *Syncthing) dispatchEvent(event events.Event) {
	s.eventSubscription.downstreamLocker.RLock()
	defer s.eventSubscription.downstreamLocker.RUnlock()
	for _, subscriber := range s.eventSubscription.subscribers[event.Type] {
		subscriber.SafeSend(event)
	}
}
