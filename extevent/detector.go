package extevent

import (
	"SyncthingHook/safechan"
	"SyncthingHook/stclient"
	"github.com/syncthing/syncthing/lib/events"
	"sync"
)

type upstream = <-chan events.Event
type subscriber = <-chan Event
type downstreamDispatcher = func(event *Event)

type EventDetector interface {
	Subscribe() subscriber
	Unsubscribe(subscriber)
}

type coreEventDetector interface {
	subscribeUpstream(*stclient.Syncthing) upstream
	unsubscribeUpstream(*stclient.Syncthing, upstream)
	// handleUpstream will be called in a coroutine and should extract the syncthing event from upstream
	handleUpstream(upstream, downstreamDispatcher)
}

type detector struct {
	syncthing   *stclient.Syncthing
	upstreamReg coreEventDetector

	locker      sync.RWMutex
	subscribers []*safechan.SafeChannel[Event]
	upstream    upstream
}

func newDetector(syncthing *stclient.Syncthing, upstreamReg coreEventDetector) *detector {
	return &detector{
		syncthing:   syncthing,
		upstreamReg: upstreamReg,
		subscribers: make([]*safechan.SafeChannel[Event], 0),
		upstream:    nil,
	}
}

func (d *detector) Subscribe() <-chan Event {
	d.locker.Lock()
	defer d.locker.Unlock()

	newSubscriber := safechan.NewSafeChannel[Event]()
	d.subscribers = append(d.subscribers, newSubscriber)
	if d.upstream == nil {
		d.upstream = d.upstreamReg.subscribeUpstream(d.syncthing)
		go d.upstreamReg.handleUpstream(d.upstream, d.dispatchEvent)
	}
	return newSubscriber.OutCh()
}

func (d *detector) Unsubscribe(ch subscriber) {
	d.locker.Lock()
	defer d.locker.Unlock()

	for i, subscriber := range d.subscribers {
		if subscriber.OutCh() == ch {
			subscriber.SafeClose()
			d.subscribers = append(d.subscribers[:i], d.subscribers[i+1:]...)
			break
		}
	}
	if len(d.subscribers) == 0 {
		d.upstreamReg.unsubscribeUpstream(d.syncthing, d.upstream)
		d.upstream = nil
	}
}

func (d *detector) dispatchEvent(newEvent *Event) {
	d.locker.RLock()
	defer d.locker.RUnlock()
	for _, sub := range d.subscribers {
		sub.SafeSend(*newEvent)
	}
}
