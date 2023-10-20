package sthook

import (
	"SyncthingHook/extevent"
	"SyncthingHook/safechan"
	"SyncthingHook/stclient"
	"github.com/syncthing/syncthing/lib/events"
	"go.uber.org/zap"
	"sync"
)

type stEventBusUpstream = <-chan events.Event

// downstream stSubscribers type
type stEventBusSubscriber = *safechan.SafeChannel[extevent.Event]
type stEventBusExposedSubscriber = <-chan extevent.Event

// stEventBus can subscribe Syncthing events and dispatch them to dowstream after converting to extevent.Event.
// stEventBus only register one upstream subscriber for each event type to avoid too many type converting coroutines.
type stEventBus struct {
	syncthing *stclient.Syncthing

	locker sync.RWMutex
	// stSubscribers records all downstream channels whose upstream is syncthing. So that we only need one upstream
	// subscription for each event type.
	stSubscribers map[events.EventType][]stEventBusSubscriber
	// stUpstreams records all upstream channels so that we can remove them if there are no downstream subscribers
	// of corresponding event type.
	stUpstreams map[events.EventType]stEventBusUpstream
	// stMapping record full-featured downstream channel so that we can close them. For security reason, users should
	// not close the output channel themselves.
	stMapping map[stEventBusExposedSubscriber]stEventBusSubscriber

	logger *zap.SugaredLogger
}

func newStEventBus(syncthing *stclient.Syncthing, logger *zap.SugaredLogger) *stEventBus {
	return &stEventBus{
		syncthing:     syncthing,
		stSubscribers: make(map[events.EventType][]stEventBusSubscriber),
		stUpstreams:   make(map[events.EventType]stEventBusUpstream),
		stMapping:     make(map[stEventBusExposedSubscriber]stEventBusSubscriber),
		logger:        logger,
	}
}

// unsubscribe close the target channel and remove the corresponding stUpstream listener if there are no more stSubscribers.
func (b *stEventBus) unsubscribe(eventCh stEventBusExposedSubscriber) {
	b.locker.Lock()
	defer b.locker.Unlock()

	realCh, ex := b.stMapping[eventCh]
	if !ex {
		b.logger.Info("given subscribers not exist, ignore unsubscribe request")
		return // already unsubscribed
	}
	// remove the stSubscribers from list
	var eventType events.EventType
Loop:
	for et, channels := range b.stSubscribers {
		for i, ch := range channels {
			if ch == realCh {
				eventType = et
				newSubscribers := append(channels[:i], channels[i+1:]...)
				if len(newSubscribers) > 0 {
					b.stSubscribers[et] = newSubscribers
				} else {
					b.logger.Debugw("subscriber list is empty, delete this type item", "eventType", et)
					delete(b.stSubscribers, et)
				}
				break Loop // safe, because one downstream subscriber can only listen to just one event type
			}
		}
	}
	// close subscriber channel
	realCh.SafeClose()
	// delete stMapping between channel returned and real channel
	delete(b.stMapping, eventCh)
	b.logger.Debugw("downstream unsubscribed", "remaining", len(b.stMapping))
	// unsubscribe stUpstream subscriber if no long needed
	if _, ex := b.stSubscribers[eventType]; !ex {
		b.logger.Debugw("unsubscribe upstream event because there's no downstream subscriber any more", "eventType", eventType)
		b.syncthing.UnsubscribeEvent(b.stUpstreams[eventType])
		delete(b.stUpstreams, eventType)
	}
}

func (b *stEventBus) subscribe(eventType events.EventType) stEventBusExposedSubscriber {
	b.locker.Lock()
	defer b.locker.Unlock()

	//  add downstream stSubscribers
	newSubscriber := safechan.NewSafeChannel[extevent.Event]()
	out := newSubscriber.OutCh()
	b.stMapping[out] = newSubscriber
	b.logger.Debugw("add subscriber to list", "eventType", eventType)
	if channels, ex := b.stSubscribers[eventType]; ex {
		b.stSubscribers[eventType] = append(channels, newSubscriber)
	} else {
		b.stSubscribers[eventType] = []*safechan.SafeChannel[extevent.Event]{newSubscriber}
	}

	if _, ex := b.stUpstreams[eventType]; ex {
		return out // already registered stUpstream listener
	}
	// register stUpstream listener
	b.logger.Debugw("register upstream event listener", "eventType", eventType)
	upstreamListener := b.syncthing.SubscribeEvent([]events.EventType{eventType}, 0)
	b.stUpstreams[eventType] = upstreamListener
	go func() {
		b.logger.Debug("start consuming upstream subscriber events")
		for event := range upstreamListener {
			// dispatch event
			b.locker.RLock()
			for _, subscribe := range b.stSubscribers[event.Type] {
				subscribe.SafeSend(extevent.NewEventFromStEvent(event))
			}
			b.locker.RUnlock()
		}
	}()
	return out
}
