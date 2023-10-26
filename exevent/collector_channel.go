package exevent

import (
	"SyncthingHook/domain"
)

// applyHandlers listens given event channel and pass the event to the handler chain built by handlers.
// Adds an ChannelCollector to transform processed events into another channel and returns it. The returned channel
// will be closed automatically once the given channel is closed. All handlers will also be destroyed at the same time.
//
// If there is no handler, the given channel is returned directly without any processing.
func applyHandlers(eventCh domain.EventCh, handlers ...domain.EventHandler) domain.EventCh {
	if len(handlers) == 0 {
		return eventCh
	}
	channelCollector := NewChannelCollector()
	chain := domain.BuildEventHandlerChain(append(handlers, channelCollector)...)
	go func() {
		for event := range eventCh {
			chain.Handle(event)
		}
		domain.DestroyHandlerChain(chain)
	}()
	return channelCollector.GetCh()
}

// ChannelCollector forwards all events to a channel and call next handler.
// The out channel will be closed once this collector is destroyed.
type ChannelCollector struct {
	baseHandler
	ch chan *domain.Event
}

var _ domain.EventHandler = (*ChannelCollector)(nil)

func NewChannelCollector() *ChannelCollector {
	return &ChannelCollector{
		ch: make(chan *domain.Event),
	}
}

func (c *ChannelCollector) GetCh() domain.EventCh {
	return c.ch
}

func (c *ChannelCollector) Handle(event *domain.Event) {
	if c.ch != nil {
		c.ch <- event
	}
	c.callNext(event)
}

func (c *ChannelCollector) Destroy() {
	if c.ch != nil {
		close(c.ch)
	}
	c.ch = nil
}
